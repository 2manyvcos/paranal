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

func (r Service) TableType() Table[Service] { return nil }

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

type ServiceTableLeftJoinFavorites struct{ JoinedTableHeader }

type ServiceWithFavorite struct {
	Service
	IsFavorite bool
}

type serviceWithFavoriteRecord struct {
	service  Service
	userName *string
}

var ServicesLeftJoinFavorites JoinedTable[ServiceWithFavorite] = ServiceTableLeftJoinFavorites{
	JoinedTableHeader{
		Type:        JoinTypeLeft,
		LeftName:    Services.TableName(),
		LeftFields:  Services.RecordFieldNames(),
		RightName:   Favorites.TableName(),
		RightFields: serviceFavoriteRightFields,
	},
}

func (t ServiceTableLeftJoinFavorites) TableCorrelations() [][2]string {
	return ServiceFavoriteCorrelations
}

var ServiceFavoriteCorrelations = [][2]string{{"id", "serviceID"}}

func (t ServiceTableLeftJoinFavorites) NewJoinedRecord() JoinedTableRecord[ServiceWithFavorite] {
	return new(serviceWithFavoriteRecord)
}

func (r *serviceWithFavoriteRecord) LeftRecordFields() []any {
	return r.service.RecordFields()
}

var serviceFavoriteRightFields = []string{"userName"}

func (r *serviceWithFavoriteRecord) RightRecordFields() []any {
	return []any{&r.userName}
}

func (r *serviceWithFavoriteRecord) Dataset() ServiceWithFavorite {
	return ServiceWithFavorite{
		Service:    r.service,
		IsFavorite: r.userName != nil,
	}
}
