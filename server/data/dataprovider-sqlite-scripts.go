package data

import (
	"fmt"

	_ "modernc.org/sqlite"
)

func (p *sqliteImpl) setupScripts() error {
	_, err := p.DB.Exec("CREATE TABLE IF NOT EXISTS scripts (id INTEGER PRIMARY KEY, name TEXT, schedule TEXT, source TEXT, serviceID INTEGER)")
	if err != nil {
		return fmt.Errorf("creating table \"scripts\" failed - %s", err)
	}
	return nil
}

func (p *sqliteImpl) ListScripts() ([]Script, error) {
	return sqliteListJoinedDatasets(p, Scripts)
}

func (p *sqliteImpl) CreateScript(record Script, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, Scripts, record, ScriptIDIDs)
	}
	return sqliteCreateDataset(p, Scripts, record)
}

func (p *sqliteImpl) ListScriptsByService(serviceID string) ([]Script, error) {
	return sqliteGetJoinedDatasets(p, Scripts, ScriptServiceID(serviceID))
}

func (p *sqliteImpl) GetScriptByService(serviceID string, id int) (result Script, err error) {
	return sqliteGetJoinedDataset(p, Scripts, ScriptIDAndServiceID{ID: id, ServiceID: serviceID})
}

func (p *sqliteImpl) UpdateScriptByService(record Script) error {
	return sqliteUpdateDataset(p, Scripts, ScriptIDAndServiceID{ID: record.ID, ServiceID: record.ServiceID}, record)
}

func (p *sqliteImpl) DeleteScriptByService(serviceID string, id int) error {
	return sqliteDeleteDataset(p, Scripts, ScriptIDAndServiceID{ID: id, ServiceID: serviceID})
}
