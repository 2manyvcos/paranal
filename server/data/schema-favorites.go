package data

import "fmt"

type FavoriteTable struct{ TableHeader }

var Favorites Table[Favorite] = FavoriteTable{
	TableHeader{
		Name:   "favorites",
		Fields: favoriteFields,
	},
}

func (t FavoriteTable) NewRecord() TableRecord[Favorite] {
	return new(Favorite)
}

type Favorite struct {
	UserName  string
	ServiceID string
}

var favoriteFields = []string{"userName", "serviceID"}

func (r *Favorite) RecordFields() []any {
	return []any{&r.UserName, &r.ServiceID}
}

func (r *Favorite) Dataset() Favorite {
	return *r
}

var _ Insertable[Favorite] = Favorite{}

func (r Favorite) TableType() Table[Favorite] { return nil }

func (r Favorite) Valid() error {
	if r.UserName == "" {
		return fmt.Errorf("invalid user name")
	}
	if r.ServiceID == "" {
		return fmt.Errorf("invalid service ID")
	}
	return nil
}

func (r Favorite) InsertableNames() []string {
	return favoriteInsertables
}

var favoriteInsertables = []string{"userName", "serviceID"}

func (r Favorite) Insertables() []any {
	return []any{r.UserName, r.ServiceID}
}
