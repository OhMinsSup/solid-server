package server

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"net/http"
	"solid-server/app"
	"solid-server/services/config"
	"solid-server/services/store"
	"solid-server/services/store/sqlstore"
	"solid-server/web"
	"sync"
)

type Server struct {
	config    *config.Configuration
	webServer *web.Server
	store     store.Store
	logger    *mlog.Logger

	servicesStartStopMutex sync.Mutex

	localRouter     *mux.Router
	localModeServer *http.Server

	app *app.App
}

func New(params Params) (*Server, error) {
	if err := params.CheckValid(); err != nil {
		return nil, err
	}

	appServices := app.Services{
		Store:  params.DBStore,
		Logger: params.Logger,
	}
	app := app.New(params.Cfg, appServices)

	// server
	webServer := web.NewServer(params.Cfg.WebPath, params.Cfg.ServerRoot, params.Cfg.Port,
		params.Cfg.UseSSL, params.Cfg.LocalOnly, params.Logger)

	server := Server{
		config:    params.Cfg,
		webServer: webServer,
		store:     params.DBStore,
		logger:    params.Logger,
		app:       app,
	}

	return &server, nil
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
		TablePrefix:      config.DBTablePrefix,
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

func (s *Server) Start() error {
	s.logger.Info("Server.Start")

	s.webServer.Start()

	s.servicesStartStopMutex.Lock()
	defer s.servicesStartStopMutex.Unlock()

	return nil
}

func (s *Server) Shutdown() error {
	if err := s.webServer.Shutdown(); err != nil {
		return err
	}

	s.servicesStartStopMutex.Lock()
	defer s.servicesStartStopMutex.Unlock()

	//s.app.Shutdown()

	defer s.logger.Info("Server.Shutdown")

	//return s.store.Shutdown()
	return nil
}

func (s *Server) Config() *config.Configuration {
	return s.config
}

func (s *Server) Logger() *mlog.Logger {
	return s.logger
}
