package sqlite

import (
	"fmt"
	"strings"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec("CREATE TABLE IF NOT EXISTS favorites (userName TEXT, serviceID TEXT, PRIMARY KEY (userName, serviceID))")
		if err != nil {
			return fmt.Errorf("creating table \"favorites\" failed - %s", err)
		}
		return nil
	})
}

func FavoriteQuery(query *schema.FavoriteQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.FavoriteQuery)
	}
	var conditions []string
	if query.UserName != nil {
		conditions = append(conditions, "userName = ?")
		placeholders = append(placeholders, *query.UserName)
	}
	if query.ServiceID != nil {
		conditions = append(conditions, "serviceID = ?")
		placeholders = append(placeholders, *query.ServiceID)
	}
	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	clause = strings.Join(conditions, " AND ")
	return
}

func (i *impl) CreateFavorite(record schema.Favorite) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO favorites (userName, serviceID)
      VALUES (?, ?)
    `,
		&record.UserName, &record.ServiceID,
	)
	return requireNoConflict(err)
}

func (i *impl) DeleteFavorite(query schema.FavoriteQuery) error {
	where, wherePlaceholders := FavoriteQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM favorites
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
