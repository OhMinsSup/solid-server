package sqlstore

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/mattermost/mattermost-plugin-api/cluster"
	"github.com/mattermost/morph"
	"text/template"

	drivers "github.com/mattermost/morph/drivers"
	mysql "github.com/mattermost/morph/drivers/mysql"
	postgres "github.com/mattermost/morph/drivers/postgres"
	sqlite "github.com/mattermost/morph/drivers/sqlite"
	embedded "github.com/mattermost/morph/sources/embedded"

	"solid-server/model"
)

//go:embed migrations
var assets embed.FS

const (
	uniqueIDsMigrationRequiredVersion      = 14
	teamsAndBoardsMigrationRequiredVersion = 17

	teamLessBoardsMigrationKey = "TeamLessBoardsMigrationComplete"

	tempSchemaMigrationTableName = "temp_schema_migration"
)

var errChannelCreatorNotInTeam = errors.New("channel creator not found in user teams")

// MySQL의 마이그레이션은 multiStatements 플래그로 실행해야 합니다.
// 활성화되어 있으므로 이 메서드는 새 연결을 생성하여
// 활성화됨.
func (s *SQLStore) getMigrationConnection() (*sql.DB, error) {
	connectionString := s.connectionString
	// TODO: MySQL Connection을 생성

	// DB 연결을 생성하고 연결을 열기
	db, err := sql.Open(s.dbType, connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *SQLStore) Migrate() error {
	var driver drivers.Driver
	var err error

	migrationConfig := drivers.Config{
		StatementTimeoutInSecs: 1000000,
		MigrationsTable: fmt.Sprintf("%sschema_migrations", s.tablePrefix),
	}

	if s.dbType == model.SqliteDBType {
		driver, err = sqlite.WithInstance(s.db, &sqlite.Config{Config: migrationConfig})
		if err != nil {
			return err
		}
	}

	var db *sql.DB
	if s.dbType != model.SqliteDBType {
		db, err = s.getMigrationConnection()
		if err != nil {
			return err
		}

		defer db.Close()
	}

	if s.dbType == model.PostgresDBType {
		driver, err = postgres.WithInstance(db, &postgres.Config{Config: migrationConfig})
		if err != nil {
			return err
		}
	}

	if s.dbType == model.MysqlDBType {
		driver, err = mysql.WithInstance(db, &mysql.Config{Config: migrationConfig})
		if err != nil {
			return err
		}
	}

	assetsList, err := assets.ReadDir("migrations")
	if err != nil {
		return err
	}

	assetNamesForDriver := make([]string, len(assetsList))
	for i, dirEntry := range assetsList {
		assetNamesForDriver[i] = dirEntry.Name()
	}

	params := map[string]interface{}{
		"prefix":   s.tablePrefix,
		"postgres": s.dbType == model.PostgresDBType,
		"sqlite":   s.dbType == model.SqliteDBType,
		"mysql":    s.dbType == model.MysqlDBType,
		"plugin":   s.isPlugin,
	}

	migrationAssets := &embedded.AssetSource{
		Names: assetNamesForDriver,
		AssetFunc: func(name string) ([]byte, error) {
			asset, mErr := assets.ReadFile("migrations/" + name)
			if mErr != nil {
				return nil, mErr
			}

			tmpl, pErr := template.New("sql").Parse(string(asset))
			if pErr != nil {
				return nil, pErr
			}
			buffer := bytes.NewBufferString("")

			err = tmpl.Execute(buffer, params)
			if err != nil {
				return nil, err
			}

			return buffer.Bytes(), nil
		},
	}

	src, err := embedded.WithInstance(migrationAssets)
	if err != nil {
		return err
	}


	opts := []morph.EngineOption{
		morph.WithLock("mm-lock-key"),
	}

	if s.dbType == model.SqliteDBType {
		opts = opts[:0] // sqlite driver does not support locking, it doesn't need to anyway.
	}

	engine, err := morph.New(context.Background(), driver, src, opts...)
	if err != nil {
		return err
	}
	defer engine.Close()

	var mutex *cluster.Mutex
	if s.isPlugin {
		var mutexErr error
		mutex, mutexErr = s.NewMutexFn("Boards_dbMutex")
		if mutexErr != nil {
			return fmt.Errorf("error creating database mutex: %w", mutexErr)
		}
	}

	if s.isPlugin {
		s.logger.Debug("Acquiring cluster lock for Unique IDs migration")
		mutex.Lock()
	}

	return engine.ApplyAll()
}
