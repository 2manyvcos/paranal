package data

import "fmt"

type ServiceTable struct{ TableHeader }

var Services Table[Service] = ServiceTable{
	TableHeader{
		Name:   "services",
		Fields: serviceFields,
	},
}

func (t ServiceTable) NewRecord() TableRecord[Service] {
	return new(Service)
}

type Service struct {
	ID          string
	Name        string
	Description string
	Logo        string
	URL         string
}

var serviceFields = []string{"id", "name", "description", "logo", "url"}

func (r *Service) RecordFields() []any {
	return []any{&r.ID, &r.Name, &r.Description, &r.Logo, &r.URL}
}

func (r *Service) Dataset() Service {
	return *r
}

var _ Upsertable[Service] = Service{}

func (r Service) Table() Table[Service] {
	return Services
}

func (r Service) Valid() error {
	if r.ID == "" {
		return fmt.Errorf("invalid ID")
	}
	if r.Name == "" {
		return fmt.Errorf("invalid name")
	}
	return nil
}

func (r Service) InsertableNames() []string {
	return serviceInsertables
}

var serviceInsertables = []string{"id", "name", "description", "logo", "url"}

func (r Service) Insertables() []any {
	return []any{r.ID, r.Name, r.Description, r.Logo, r.URL}
}

func (r Service) UpdatableNames() []string {
	return serviceUpdatables
}

var serviceUpdatables = []string{"name", "description", "logo", "url"}

func (r Service) Updatables() []any {
	return []any{r.Name, r.Description, r.Logo, r.URL}
}

type ServiceID string

var _ Identifier[Service] = ServiceID("")

func (i ServiceID) Table() Table[Service] {
	return Services
}

func (i ServiceID) Valid() error {
	if i == "" {
		return fmt.Errorf("invalid ID")
	}
	return nil
}

func (i ServiceID) IDNames() []string {
	return ServiceIDIDs
}

var ServiceIDIDs = []string{"id"}

func (i ServiceID) IDs() []any {
	return []any{string(i)}
}
