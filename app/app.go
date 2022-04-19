package app

import (
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/services/store"
	"time"
)

const (
	blockChangeNotifierQueueSize       = 100
	blockChangeNotifierPoolSize        = 10
	blockChangeNotifierShutdownTimeout = time.Second * 10
)

type App struct {
	Store  store.Store
	logger *mlog.Logger
	auth   interface{}
}
