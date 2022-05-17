package sqlstore

import (
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/model"
	"solid-server/utils"
)

var (
	teamFields = []string{
		"id",
		"signup_token",
		"COALESCE(settings, '{}')",
		"modified_by",
		"update_at",
	}
)

func (s *SQLStore) upsertTeamSignupToken(db sq.BaseRunner, team model.Team) error {
	now := utils.GetMillis()

	query := s.getQueryBuilder(db).
		Insert(s.tablePrefix+"teams").
		Columns(
			"id",
			"signup_token",
			"modified_by",
			"update_at",
		).
		Values(team.ID, team.SignupToken, team.ModifiedBy, now)

	if s.dbType == model.MysqlDBType {
		query = query.Suffix("ON DUPLICATE KEY UPDATE signup_token = ?, modified_by = ?, update_at = ?",
			team.SignupToken, team.ModifiedBy, now)
	} else {
		query = query.Suffix(
			`ON CONFLICT (id)
			 DO UPDATE SET signup_token = EXCLUDED.signup_token, modified_by = EXCLUDED.modified_by, update_at = EXCLUDED.update_at`,
		)
	}

	_, err := query.Exec()
	return err
}

func (s *SQLStore) getTeam(db sq.BaseRunner, id string) (*model.Team, error) {
	var settingsJSON string

	query := s.getQueryBuilder(db).
		Select(
			teamFields...,
		).
		From(s.tablePrefix + "teams").
		Where(sq.Eq{"id": id})
	row := query.QueryRow()
	team := model.Team{}

	err := row.Scan(
		&team.ID,
		&team.SignupToken,
		&settingsJSON,
		&team.ModifiedBy,
		&team.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(settingsJSON), &team.Settings)
	if err != nil {
		s.logger.Error(`ERROR GetTeam settings json.Unmarshal`, mlog.Err(err))
		return nil, err
	}

	return &team, nil
}
