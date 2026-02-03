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

func (p *sqliteImpl) ListServices() ([]Service, error) {
	return sqliteListDatasets(p, Services)
}

func (p *sqliteImpl) GetService(id string) (result Service, err error) {
	return sqliteGetDataset(p, Services, ServiceID(id))
}

func (p *sqliteImpl) CreateService(record Service, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, Services, record, ServiceIDIDs)
	}
	return sqliteCreateDataset(p, Services, record)
}

func (p *sqliteImpl) UpdateService(record Service) error {
	return sqliteUpdateDataset(p, Services, ServiceID(record.ID), record)
}

func (p *sqliteImpl) DeleteService(id string) error {
	return sqliteDeleteDataset(p, Services, ServiceID(id))
}
