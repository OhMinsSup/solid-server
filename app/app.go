package app

import (
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/services/auth"
	"solid-server/services/config"
	"solid-server/services/store"
	"time"
)

const (
	blockChangeNotifierQueueSize       = 100
	blockChangeNotifierPoolSize        = 10
	blockChangeNotifierShutdownTimeout = time.Second * 10
)

type Services struct {
	Auth   *auth.Auth
	Store  store.Store
	Logger *mlog.Logger
}

type App struct {
	config *config.Configuration
	store  store.Store
	logger *mlog.Logger
	auth   *auth.Auth
}

func (a *App) SetConfig(config *config.Configuration) {
	a.config = config
}

func (a *App) GetConfig() *config.Configuration {
	return a.config
}

func (a *App) Shutdown() {
	a.logger.Info("Shutting down")
	// TODO: actually shutdown
}

func New(config *config.Configuration, services Services) *App {
	app := &App{
		config: config,
		store:  services.Store,
		logger: services.Logger,
		auth:   services.Auth,
	}
	return app
}
