package sqlstore

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"solid-server/model"
	"solid-server/utils"
)

func (s *SQLStore) getFindOrCreate(db sq.BaseRunner, names []string) error {
	categoriesQuery := s.getQueryBuilder(db).
		Select("id", "name").
		From(s.tablePrefix + "categories").
		Where(sq.Eq{"name": names})

	rows, err := categoriesQuery.Query()
	if err != nil {
		s.logger.Error(`getFindOrCreate ERROR`, mlog.Err(err))
		return err
	}
	defer s.CloseRows(rows)

	categories, err := s.categoriesFromRows(rows)
	if err != nil {
		s.logger.Error(`getFindOrCreate Row Scan ERROR`, mlog.Err(err))
		return err
	}

	fmt.Println("ca", categories)
	fmt.Println("names", names)

	if len(categories) == 0 {
		insertBuilder := s.getQueryBuilder(db).Insert(s.tablePrefix+"categories").Columns("id", "name")
		for _, name := range names {
			insertBuilder = insertBuilder.Values(name, utils.NewID(utils.IDTypeCategories))
		}

		_, err = insertBuilder.Exec()
		if err != nil {
			s.logger.Error(`getFindOrCreate categories table insert ERROR`, mlog.Err(err))
			return err
		}
	}

	return nil
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
