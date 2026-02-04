package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupFavorites() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS favorites (id INTEGER PRIMARY KEY, userName TEXT, serviceID TEXT, UNIQUE(userName, serviceID))")
	if err != nil {
		return fmt.Errorf("creating table \"favorites\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) CreateFavorite(record Favorite) error {
	return sqliteCreateDataset(p, Favorites, record)
}

func (p *sqliteImpl) DeleteFavorite(userName string, serviceID string) error {
	return sqliteDeleteDataset(p, Favorites, &Conditions{Conditions: []Conditions{{Condition: &Condition{Field: "userName", Value: userName}}, {Condition: &Condition{Field: "serviceID", Value: serviceID}}}})
}
