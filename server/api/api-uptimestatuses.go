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

type UptimeStatus struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Unhealthy bool   `json:"unhealthy"`
}

type ServiceUptimeStatus struct {
	UptimeStatus
	ServiceID string `json:"serviceID"`
}

func GetUptimeStatuses(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil || authorizedUser.Name == "" {
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
	if v, ok := utils.LoadQueryValue(q, "name"); ok {
		nameQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "status"); ok {
		t, ok := schema.ServiceUptimeStatusCodes[v]
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

	uptimeStatuses := make([][]ServiceUptimeStatus, len(services))
	for i, service := range services {
		records := app.ListServiceUptimeStatuses(service.ID)
		uptimeStatuses[i] = make([]ServiceUptimeStatus, 0, len(records))
		for _, uptimeStatus := range records {
			if nameQuery != nil && uptimeStatus.Name != *nameQuery {
				continue
			}
			if statusQuery != nil && uptimeStatus.Status != *statusQuery {
				continue
			}
			if unhealthyQuery != nil && uptimeStatus.Unhealthy != *unhealthyQuery {
				continue
			}
			uptimeStatuses[i] = append(uptimeStatuses[i], ServiceUptimeStatus{
				UptimeStatus: UptimeStatus{
					Name:      uptimeStatus.Name,
					Status:    schema.ServiceUptimeStatusNames[uptimeStatus.Status],
					Unhealthy: uptimeStatus.Unhealthy,
				},
				ServiceID: service.ID,
			})
		}
	}
	responsePayload := slices.Concat(uptimeStatuses...)
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func GetServicesByIDUptimeStatuses(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil || authorizedUser.Name == "" {
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
		t, ok := schema.ServiceUptimeStatusCodes[v]
		if !ok {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		statusQuery = &t
	}
	if v, ok := utils.LoadQueryBool(q, "unhealthy"); ok {
		unhealthyQuery = &v
	}

	records := app.ListServiceUptimeStatuses(serviceID)
	responsePayload := make([]UptimeStatus, 0, len(records))
	for _, uptimeStatus := range records {
		if nameQuery != nil && uptimeStatus.Name != *nameQuery {
			continue
		}
		if statusQuery != nil && uptimeStatus.Status != *statusQuery {
			continue
		}
		if unhealthyQuery != nil && uptimeStatus.Unhealthy != *unhealthyQuery {
			continue
		}
		responsePayload = append(responsePayload, UptimeStatus{
			Name:      uptimeStatus.Name,
			Status:    schema.ServiceUptimeStatusNames[uptimeStatus.Status],
			Unhealthy: uptimeStatus.Unhealthy,
		})
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}
