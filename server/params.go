package server

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/services/config"
	"solid-server/services/store"
)

type Params struct {
	Cfg             *config.Configuration
	DBStore         store.Store
	Logger          *mlog.Logger
	ServerID        string
}

func (p Params) CheckValid() error {
	if p.Cfg == nil {
		return ErrServerParam{name: "Cfg", issue: "cannot be nil"}
	}

	if p.DBStore == nil {
		return ErrServerParam{name: "DbStore", issue: "cannot be nil"}
	}

	if p.Logger == nil {
		return ErrServerParam{name: "Logger", issue: "cannot be nil"}
	}

	return nil
}

type ErrServerParam struct {
	name  string
	issue string
}

func (e ErrServerParam) Error() string {
	return fmt.Sprintf("invalid server params: %s %s", e.name, e.issue)
}
