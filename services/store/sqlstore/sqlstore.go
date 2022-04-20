package sqlstore

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-plugin-api/cluster"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"solid-server/model"

	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

// SQLStore SQL 데이터베이스
type SQLStore struct {
	db               *sql.DB
	dbType           string
	tablePrefix      string
	connectionString string
	isPlugin         bool
	logger           *mlog.Logger
	NewMutexFn       MutexFactory
	pluginAPI        *plugin.API
}

// MutexFactory 플러그인 모드에서 스토어에서 생성하는 데 사용됩니다. cluster mutex.
type MutexFactory func(name string) (*cluster.Mutex, error)

// New 저장소의 새 SQL 구현을 만듭니다.
func New(params Params) (*SQLStore, error) {
	if err := params.CheckValid(); err != nil {
		return nil, err
	}

	params.Logger.Info("connectDatabase", mlog.String("dbType", params.DBType))

	store := &SQLStore{
		// TODO: add replica DB support too.
		db:               params.DB,
		dbType:           params.DBType,
		tablePrefix:      params.TablePrefix,
		connectionString: params.ConnectionString,
		logger:           params.Logger,
		isPlugin:         params.IsPlugin,
		NewMutexFn:       params.NewMutexFn,
		pluginAPI:        params.PluginAPI,
	}

	err := store.Migrate()
	if err != nil {
		params.Logger.Error(`Table creation / migration failed`, mlog.Err(err))

		return nil, err
	}

	return store, nil
}

// Shutdown 스토어와의 연결을 닫습니다.
func (s *SQLStore) Shutdown() error {
	return s.db.Close()
}

// DBHandle raw sql.DB 핸들을 반환합니다.
// Mattermostauthlayer가 자체 실행에 사용합니다.
// raw SQL 쿼리.
func (s *SQLStore) DBHandle() *sql.DB {
	return s.db
}

// DBType 저장소에 사용된 DB 드라이버를 반환합니다.
func (s *SQLStore) DBType() string {
	return s.dbType
}

func (s *SQLStore) getQueryBuilder(db sq.BaseRunner) sq.StatementBuilderType {
	builder := sq.StatementBuilder
	if s.dbType == model.PostgresDBType || s.dbType == model.SqliteDBType {
		builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	}
	return builder.RunWith(db)
}

func (s *SQLStore) escapeField(fieldName string) string { //nolint:unparam
	if s.dbType == model.MysqlDBType {
		return "`" + fieldName + "`"
	}
	if s.dbType == model.PostgresDBType || s.dbType == model.SqliteDBType {
		return "\"" + fieldName + "\""
	}
	return fieldName
}
