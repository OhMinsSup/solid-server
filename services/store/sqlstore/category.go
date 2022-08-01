package sqlstore

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/model"
	"solid-server/utils"
)

func (s *SQLStore) findCategories(db sq.BaseRunner, names []string) ([]*model.Category, error) {
	categoriesQuery := s.getQueryBuilder(db).
		Select("id", "name").
		From(s.tablePrefix + "categories").
		Where(sq.Eq{"name": names})

	rows, err := categoriesQuery.Query()

	if err != nil {
		s.logger.Error(`findCategories ERROR`, mlog.Err(err))
		return nil, err
	}
	defer s.CloseRows(rows)

	categories, err := s.categoriesFromRows(rows)
	if err != nil {
		s.logger.Error(`findCategories Row Scan ERROR`, mlog.Err(err))
		return nil, err
	}

	return categories, nil
}

func (s *SQLStore) getFindOrCreate(db sq.BaseRunner, names []string) ([]*model.Category, error) {
	now := utils.GetMillis()
	categories, err := s.findCategories(db, names)
	if err != nil {
		s.logger.Error(`getFindOrCreate findCategories`, mlog.Err(err))
		return nil, err
	}

	if len(categories) == 0 {
		insertBuilder := s.getQueryBuilder(db).Insert(s.tablePrefix+"categories").Columns("id", "name", "create_at", "update_at")
		for _, name := range names {
			insertBuilder = insertBuilder.Values(utils.NewID(utils.IDTypeCategories), name, now, now, 0)
		}

		_, err = insertBuilder.Exec()
		if err != nil {
			s.logger.Error(`getFindOrCreate categories table insert ERROR`, mlog.Err(err))
			return nil, err
		}
	} else {
		var others []string
		for _, name := range names {
			if s.containsByCategory(categories, name) {
				// 등록이 된 경우
			} else {
				// 등록이 안된 경우
				others = append(others, name)
			}
		}

		if len(others) != 0 {
			insertBuilder := s.getQueryBuilder(db).Insert(s.tablePrefix+"categories").Columns("id", "name", "create_at", "update_at")
			for _, other := range others {
				insertBuilder = insertBuilder.Values(utils.NewID(utils.IDTypeCategories), other, now, now, 0)
			}

			_, err = insertBuilder.Exec()
			if err != nil {
				s.logger.Error(`getFindOrCreate categories table insert ERROR`, mlog.Err(err))
				return nil, err
			}
		}
	}

	categories, err = s.findCategories(db, names)
	if err != nil {
		s.logger.Error(`getFindOrCreate findCategories`, mlog.Err(err))
		return nil, err
	}

	return categories, nil
}

// ContainsByString 는 string 문자열중에 일치하는 값이 존재하는지 체크
func (s *SQLStore) containsByCategory(categories []*model.Category, str string) bool {
	for _, category := range categories {
		if category.Name == str {
			return true
		}
	}
	return false
}

func (s *SQLStore) categoriesFromRows(rows *sql.Rows) ([]*model.Category, error) {
	categories := []*model.Category{}

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}
