package sqlstore

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/mattermost/morph/drivers"
	"github.com/mattermost/morph/drivers/mysql"
	"github.com/mattermost/morph/drivers/postgres"
	"github.com/mattermost/morph/drivers/sqlite"
	"log"
	"solid-server/model"
)

//go:embed migrations
var assets embed.FS

const (

)

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

	log.Println("Migrating database schema to version", params, driver)

	return nil
}
