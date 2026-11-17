package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/server/db"
)

func GetRegistriesByNameStates(res http.ResponseWriter, req *http.Request) {
	app := GetApp(req)

	registryName := req.PathValue("registryName")
	if registryName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	records, err := app.Store.ListStates(registryName)
	if err != nil {
		log.Printf("loading records: %v\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responsePayload := make([]struct {
		Registry string         `json:"registry"`
		ID       string         `json:"id"`
		Data     map[string]any `json:"data"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].Registry = record.Registry
		responsePayload[i].ID = record.ID
		responsePayload[i].Data = record.Data
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("encoding response payload: %v\n", err)
	}
}

func GetRegistriesByNameStatesByID(res http.ResponseWriter, req *http.Request) {
	app := GetApp(req)

	registryName := req.PathValue("registryName")
	stateID := req.PathValue("stateID")
	if registryName == "" || stateID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record, err := app.Store.GetState(registryName, stateID)
	if errors.Is(err, db.ErrNotFound) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("loading record: %v\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		Registry string         `json:"registry"`
		ID       string         `json:"id"`
		Data     map[string]any `json:"data"`
	}{
		Registry: record.Registry,
		ID:       record.ID,
		Data:     record.Data,
	})
	if err != nil {
		log.Printf("encoding response payload: %v\n", err)
	}
}
