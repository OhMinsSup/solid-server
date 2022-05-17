package app

import (
	"database/sql"
	"errors"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/model"
	"solid-server/utils"
)

func (a *App) GetRootTeam() (*model.Team, error) {
	teamID := "0"
	team, _ := a.store.GetTeam(teamID)
	if team == nil {
		team = &model.Team{
			ID:          teamID,
			SignupToken: utils.NewID(utils.IDTypeToken),
		}
		err := a.store.UpsertTeamSignupToken(*team)
		if err != nil {
			a.logger.Error("Unable to initialize team", mlog.Err(err))
			return nil, err
		}

		team, err = a.store.GetTeam(teamID)
		if err != nil {
			a.logger.Error("Unable to get initialized team", mlog.Err(err))
			return nil, err
		}

		a.logger.Info("initialized team")
	}
	return team, nil
}

func (a *App) GetTeam(id string) (*model.Team, error) {
	team, err := a.store.GetTeam(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (a *App) UpsertTeamSignupToken(team model.Team) error {
	return a.store.UpsertTeamSignupToken(team)
}
