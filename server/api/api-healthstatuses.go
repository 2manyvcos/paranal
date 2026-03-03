package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
)

type HealthStatus struct {
	Name      string `json:"name"`
	ServiceID string `json:"serviceID"`
	Status    string `json:"status"`
	Unhealthy bool   `json:"unhealthy"`
}

func GetHealthStatuses(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var serviceQuery schema.UserServiceQuery
	var nameQuery *string
	var statusQuery *int
	var unhealthyQuery *bool
	q := req.URL.Query()
	if v, ok := utils.LoadQueryValue(q, "serviceID"); ok {
		serviceQuery.ID = &v
	}
	if v, ok := utils.LoadQueryBool(q, "service.config.hidden"); ok {
		serviceQuery.Hidden = &v
	}
	if v, ok := utils.LoadQueryValue(q, "name"); ok {
		nameQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "status"); ok {
		t, ok := schema.ServiceHealthStatusCodes[v]
		if !ok {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		statusQuery = &t
	}
	if v, ok := utils.LoadQueryBool(q, "unhealthy"); ok {
		unhealthyQuery = &v
	}

	services, err := app.ListUserServices(authorizedUser.Name, &serviceQuery)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	healthStatuses := make([][]HealthStatus, len(services))
	for i, service := range services {
		records := app.ListServiceHealthStatuses(service.ID)
		healthStatuses[i] = make([]HealthStatus, 0, len(records))
		for _, healthStatus := range records {
			if nameQuery != nil && healthStatus.Name != *nameQuery {
				continue
			}
			if statusQuery != nil && healthStatus.Status != *statusQuery {
				continue
			}
			if unhealthyQuery != nil && healthStatus.Unhealthy != *unhealthyQuery {
				continue
			}
			healthStatuses[i] = append(healthStatuses[i], HealthStatus{
				Name:      healthStatus.Name,
				ServiceID: service.ID,
				Status:    schema.ServiceHealthStatusNames[healthStatus.Status],
				Unhealthy: healthStatus.Unhealthy,
			})
		}
	}
	responsePayload := slices.Concat(healthStatuses...)
	if responsePayload == nil {
		responsePayload = []HealthStatus{}
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func GetServicesByIDHealthStatuses(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var nameQuery *string
	var statusQuery *int
	var unhealthyQuery *bool
	q := req.URL.Query()
	if v, ok := utils.LoadQueryValue(q, "name"); ok {
		nameQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "status"); ok {
		t, ok := schema.ServiceHealthStatusCodes[v]
		if !ok {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		statusQuery = &t
	}
	if v, ok := utils.LoadQueryBool(q, "unhealthy"); ok {
		unhealthyQuery = &v
	}

	records := app.ListServiceHealthStatuses(serviceID)
	responsePayload := make([]HealthStatus, 0, len(records))
	for _, healthStatus := range records {
		if nameQuery != nil && healthStatus.Name != *nameQuery {
			continue
		}
		if statusQuery != nil && healthStatus.Status != *statusQuery {
			continue
		}
		if unhealthyQuery != nil && healthStatus.Unhealthy != *unhealthyQuery {
			continue
		}
		responsePayload = append(responsePayload, HealthStatus{
			Name:      healthStatus.Name,
			ServiceID: serviceID,
			Status:    schema.ServiceHealthStatusNames[healthStatus.Status],
			Unhealthy: healthStatus.Unhealthy,
		})
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}
