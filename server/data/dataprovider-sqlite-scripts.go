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

func (p *sqliteImpl) GetScript(name string) (result Script, err error) {
	return sqliteGetJoinedDataset(p, Scripts, ScriptName(name))
}

func (p *sqliteImpl) CreateScript(record Script, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, Scripts, record, ScriptNameIDs)
	}
	return sqliteCreateDataset(p, Scripts, record)
}

func (p *sqliteImpl) UpdateScript(record Script) error {
	return sqliteUpdateDataset(p, Scripts, ScriptName(record.Name), record)
}

func (p *sqliteImpl) DeleteScript(name string) error {
	return sqliteDeleteDataset(p, Scripts, ScriptName(name))
}
