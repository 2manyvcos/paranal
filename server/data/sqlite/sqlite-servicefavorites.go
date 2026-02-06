package sqlite

import (
	"fmt"
	"strings"

	"github.com/2manyvcos/paranal/server/data/schema"
)

func init() {
	setups = append(setups, func(i *impl) error {
		_, err := i.Exec("CREATE TABLE IF NOT EXISTS servicefavorites (userName TEXT, serviceID TEXT, PRIMARY KEY (userName, serviceID))")
		if err != nil {
			return fmt.Errorf("creating table \"servicefavorites\" failed - %s", err)
		}
		return nil
	})
}

func ServiceFavoriteQuery(query *schema.ServiceFavoriteQuery) (clause string, placeholders []any) {
	if query == nil {
		query = new(schema.ServiceFavoriteQuery)
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

func (i *impl) CreateServiceFavorite(record schema.ServiceFavorite) error {
	if err := record.Valid(); err != nil {
		return err
	}
	_, err := i.Exec(
		`
      INSERT INTO servicefavorites (userName, serviceID)
      VALUES (?, ?)
    `,
		&record.UserName, &record.ServiceID,
	)
	return requireNoConflict(err)
}

func (i *impl) DeleteServiceFavorite(query schema.ServiceFavoriteQuery) error {
	where, wherePlaceholders := ServiceFavoriteQuery(&query)
	result, err := i.Exec(
		`
      DELETE FROM servicefavorites
      WHERE `+where+`
    `,
		wherePlaceholders...,
	)
	return requireFound(result, err)
}
