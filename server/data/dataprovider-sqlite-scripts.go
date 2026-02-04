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
	return sqliteSelectJoinedDatasets(p, ScriptsInnerJoinServices, nil, nil)
}

func (p *sqliteImpl) CreateScript(record Script, updateExisting bool) error {
	if updateExisting {
		return sqliteCreateOrUpdateDataset(p, Scripts, record, []string{"id"})
	}
	return sqliteCreateDataset(p, Scripts, record)
}

func (p *sqliteImpl) ListScriptsByService(serviceID string) ([]Script, error) {
	return sqliteSelectJoinedDatasets(p, ScriptsInnerJoinServices, nil, &JoinedConditions{LeftCondition: &Condition{Field: "serviceID", Value: serviceID}})
}

func (p *sqliteImpl) GetScriptByService(id int, serviceID string) (result Script, err error) {
	return sqliteSelectJoinedDataset(p, ScriptsInnerJoinServices, nil, &JoinedConditions{Conditions: []JoinedConditions{{LeftCondition: &Condition{Field: "id", Value: id}}, {LeftCondition: &Condition{Field: "serviceID", Value: serviceID}}}})
}

func (p *sqliteImpl) UpdateScriptByService(id int, serviceID string, record Script) error {
	return sqliteUpdateDataset(p, Scripts, &Conditions{Conditions: []Conditions{{Condition: &Condition{Field: "id", Value: id}}, {Condition: &Condition{Field: "serviceID", Value: serviceID}}}}, record)
}

func (p *sqliteImpl) DeleteScriptByService(id int, serviceID string) error {
	return sqliteDeleteDataset(p, Scripts, &Conditions{Conditions: []Conditions{{Condition: &Condition{Field: "id", Value: id}}, {Condition: &Condition{Field: "serviceID", Value: serviceID}}}})
}
