package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupServices() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS services (id TEXT PRIMARY KEY, name TEXT, description TEXT, logo TEXT, url TEXT)")
	if err != nil {
		return fmt.Errorf("creating table \"services\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) GetService(id string) (result Service, err error) {
	return sqliteSelectDataset(p, Services, &Conditions{Condition: &Condition{Field: "id", Value: id}})
}

func (p *sqliteImpl) CreateService(record Service, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, Services, record, []string{"id"})
	}
	return sqliteCreateDataset(p, Services, record)
}

func (p *sqliteImpl) UpdateService(id string, record Service) error {
	return sqliteUpdateDataset(p, Services, &Conditions{Condition: &Condition{Field: "id", Value: id}}, record)
}

func (p *sqliteImpl) DeleteService(id string) error {
	return sqliteDeleteDataset(p, Services, &Conditions{Condition: &Condition{Field: "id", Value: id}})
}

func (p *sqliteImpl) ListServicesWithFavorite(userName string) ([]ServiceWithFavorite, error) {
	return sqliteSelectJoinedDatasets(p, ServicesLeftJoinFavorites, &JoinedConditions{RightCondition: &Condition{Field: "userName", Value: userName}}, nil)
}

func (p *sqliteImpl) GetServiceWithFavorite(id string, userName string) (result ServiceWithFavorite, err error) {
	return sqliteSelectJoinedDataset(p, ServicesLeftJoinFavorites, &JoinedConditions{RightCondition: &Condition{Field: "userName", Value: userName}}, &JoinedConditions{LeftCondition: &Condition{Field: "id", Value: id}})
}
