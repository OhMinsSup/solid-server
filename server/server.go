package server

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"net/http"
	"solid-server/services/config"
	"solid-server/services/store"
	"solid-server/services/store/sqlstore"
	"sync"
)

type Server struct {
	config *config.Configuration
	logger *mlog.Logger
	store  store.Store

	servicesStartStopMutex sync.Mutex

	localRouter     *mux.Router
	localModeServer *http.Server
}

func NewStore(config *config.Configuration, logger *mlog.Logger) (store.Store, error) {
	sqlDB, err := sql.Open(config.DBType, config.DBConfigString)
	if err != nil {
		logger.Error("connectDatabase failed", mlog.Err(err))
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		logger.Error(`Database Ping failed`, mlog.Err(err))
		return nil, err
	}

	storeParams := sqlstore.Params{
		DBType:           config.DBType,
		ConnectionString: config.DBConfigString,
		TablePrefix:       config.DBTablePrefix,
		Logger:           logger,
		DB:               sqlDB,
		IsPlugin:         false,
	}

	var db store.Store
	db, err = sqlstore.New(storeParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}
