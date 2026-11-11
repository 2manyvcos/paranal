package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type dashboardRoute struct {
	Host string `json:"host,omitempty"`
	Path string `json:"path"`
}

type dashboardVar struct {
	Value   any  `json:"value"`
	Exposed bool `json:"exposed"`
}

func GetDashboards(res http.ResponseWriter, req *http.Request) {
	app := GetApp(req)

	records := app.ListDashboards()

	responsePayload := make([]struct {
		Version       string                  `json:"version"`
		ID            string                  `json:"id"`
		Routes        []dashboardRoute        `json:"routes,omitempty"`
		Theme         map[string]any          `json:"theme,omitempty"`
		Vars          map[string]dashboardVar `json:"vars,omitempty"`
		DefaultLayout map[string]any          `json:"defaultLayout,omitempty"`
		Layout        map[string]any          `json:"layout,omitempty"`
		Content       []map[string]any        `json:"content,omitempty"`
		Scripts       map[string]string       `json:"scripts,omitempty"`
	}, len(records))
	for i, record := range records {
		responsePayload[i].Version = record.Version
		responsePayload[i].ID = record.ID
		responsePayload[i].Routes = make([]dashboardRoute, len(record.Routes))
		for j, route := range record.Routes {
			responsePayload[i].Routes[j].Host = route.Host
			responsePayload[i].Routes[j].Path = route.Path
		}
		responsePayload[i].Theme = record.Theme
		responsePayload[i].Vars = make(map[string]dashboardVar, len(record.Vars))
		for name, _var := range record.Vars {
			responsePayload[i].Vars[name] = dashboardVar{
				Value:   _var.Value,
				Exposed: _var.Exposed,
			}
		}
		responsePayload[i].DefaultLayout = record.DefaultLayout
		responsePayload[i].Layout = record.Layout
		responsePayload[i].Content = record.Content
		responsePayload[i].Scripts = record.Scripts
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("encoding response payload: %v\n", err)
	}
}

func GetDashboardsByID(res http.ResponseWriter, req *http.Request) {
	app := GetApp(req)

	dashboardID := req.PathValue("dashboardID")
	if dashboardID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	record := app.GetDashboard(dashboardID)
	if record == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var responsePayload struct {
		Version       string                  `json:"version"`
		ID            string                  `json:"id"`
		Routes        []dashboardRoute        `json:"routes,omitempty"`
		Theme         map[string]any          `json:"theme,omitempty"`
		Vars          map[string]dashboardVar `json:"vars,omitempty"`
		DefaultLayout map[string]any          `json:"defaultLayout,omitempty"`
		Layout        map[string]any          `json:"layout,omitempty"`
		Content       []map[string]any        `json:"content,omitempty"`
		Scripts       map[string]string       `json:"scripts,omitempty"`
	}
	responsePayload.Version = record.Version
	responsePayload.ID = record.ID
	responsePayload.Routes = make([]dashboardRoute, len(record.Routes))
	for j, route := range record.Routes {
		responsePayload.Routes[j].Host = route.Host
		responsePayload.Routes[j].Path = route.Path
	}
	responsePayload.Theme = record.Theme
	responsePayload.Vars = make(map[string]dashboardVar, len(record.Vars))
	for name, _var := range record.Vars {
		responsePayload.Vars[name] = dashboardVar{
			Value:   _var.Value,
			Exposed: _var.Exposed,
		}
	}
	responsePayload.DefaultLayout = record.DefaultLayout
	responsePayload.Layout = record.Layout
	responsePayload.Content = record.Content
	responsePayload.Scripts = record.Scripts
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("encoding response payload: %v\n", err)
	}
}
