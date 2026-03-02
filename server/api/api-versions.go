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

type Version struct {
	Name           string `json:"name"`
	ServiceID      string `json:"serviceID"`
	CurrentVersion string `json:"currentVersion"`
	CurrentCVEs    int    `json:"currentCVEs"`
	LatestVersion  string `json:"latestVersion"`
	LatestCVEs     int    `json:"latestCVEs"`
	Status         string `json:"status"`
	Outdated       bool   `json:"outdated"`
	Vulnerable     bool   `json:"vulnerable"`
}

type VersionDetails struct {
	CurrentVersionNotes    string `json:"currentVersionNotes"`
	CurrentCVEDescriptions []CVE  `json:"currentCVEDescriptions"`
	LatestVersionNotes     string `json:"latestVersionNotes"`
	LatestCVEDescriptions  []CVE  `json:"latestCVEDescriptions"`
}

type CVE struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func GetVersions(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var serviceQuery schema.UserServiceQuery
	var nameQuery *string
	var currentVersionQuery *string
	var hasCurrentCVEsQuery *bool
	var latestVersionQuery *string
	var hasLatestCVEsQuery *bool
	var statusQuery *int
	var outdatedQuery *bool
	var vulnerableQuery *bool
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
	if v, ok := utils.LoadQueryValue(q, "currentVersion"); ok {
		currentVersionQuery = &v
	}
	if v, ok := utils.LoadQueryBool(q, "hasCurrentCVEs"); ok {
		hasCurrentCVEsQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "latestVersion"); ok {
		latestVersionQuery = &v
	}
	if v, ok := utils.LoadQueryBool(q, "hasLatestCVEs"); ok {
		hasLatestCVEsQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "status"); ok {
		t, ok := schema.ServiceVersionStatusCodes[v]
		if !ok {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		statusQuery = &t
	}
	if v, ok := utils.LoadQueryBool(q, "outdated"); ok {
		outdatedQuery = &v
	}
	if v, ok := utils.LoadQueryBool(q, "vulnerable"); ok {
		vulnerableQuery = &v
	}

	services, err := app.ListUserServices(authorizedUser.Name, &serviceQuery)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	versions := make([][]Version, len(services))
	for i, service := range services {
		records := app.ListServiceVersions(service.ID)
		versions[i] = make([]Version, 0, len(records))
		for _, version := range records {
			if nameQuery != nil && version.Name != *nameQuery {
				continue
			}
			if currentVersionQuery != nil && version.CurrentVersion != *currentVersionQuery {
				continue
			}
			if hasCurrentCVEsQuery != nil && (version.CurrentCVEs > 0) != *hasCurrentCVEsQuery {
				continue
			}
			if latestVersionQuery != nil && version.LatestVersion != *latestVersionQuery {
				continue
			}
			if hasLatestCVEsQuery != nil && (version.LatestCVEs > 0) != *hasLatestCVEsQuery {
				continue
			}
			if statusQuery != nil && version.Status != *statusQuery {
				continue
			}
			if outdatedQuery != nil && version.Outdated != *outdatedQuery {
				continue
			}
			if vulnerableQuery != nil && version.Vulnerable != *vulnerableQuery {
				continue
			}
			versions[i] = append(versions[i], Version{
				Name:           version.Name,
				ServiceID:      service.ID,
				CurrentVersion: version.CurrentVersion,
				CurrentCVEs:    version.CurrentCVEs,
				LatestVersion:  version.LatestVersion,
				LatestCVEs:     version.LatestCVEs,
				Status:         schema.ServiceVersionStatusNames[version.Status],
				Outdated:       version.Outdated,
				Vulnerable:     version.Vulnerable,
			})
		}
	}
	responsePayload := slices.Concat(versions...)
	if responsePayload == nil {
		responsePayload = []Version{}
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func GetServicesByIDVersions(res http.ResponseWriter, req *http.Request) {
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
	var currentVersionQuery *string
	var hasCurrentCVEsQuery *bool
	var latestVersionQuery *string
	var hasLatestCVEsQuery *bool
	var statusQuery *int
	var outdatedQuery *bool
	var vulnerableQuery *bool
	q := req.URL.Query()
	if v, ok := utils.LoadQueryValue(q, "name"); ok {
		nameQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "currentVersion"); ok {
		currentVersionQuery = &v
	}
	if v, ok := utils.LoadQueryBool(q, "hasCurrentCVEs"); ok {
		hasCurrentCVEsQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "latestVersion"); ok {
		latestVersionQuery = &v
	}
	if v, ok := utils.LoadQueryBool(q, "hasLatestCVEs"); ok {
		hasLatestCVEsQuery = &v
	}
	if v, ok := utils.LoadQueryValue(q, "status"); ok {
		t, ok := schema.ServiceVersionStatusCodes[v]
		if !ok {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		statusQuery = &t
	}
	if v, ok := utils.LoadQueryBool(q, "outdated"); ok {
		outdatedQuery = &v
	}
	if v, ok := utils.LoadQueryBool(q, "vulnerable"); ok {
		vulnerableQuery = &v
	}

	records := app.ListServiceVersions(serviceID)
	responsePayload := make([]Version, 0, len(records))
	for _, version := range records {
		if nameQuery != nil && version.Name != *nameQuery {
			continue
		}
		if currentVersionQuery != nil && version.CurrentVersion != *currentVersionQuery {
			continue
		}
		if hasCurrentCVEsQuery != nil && (version.CurrentCVEs > 0) != *hasCurrentCVEsQuery {
			continue
		}
		if latestVersionQuery != nil && version.LatestVersion != *latestVersionQuery {
			continue
		}
		if hasLatestCVEsQuery != nil && (version.LatestCVEs > 0) != *hasLatestCVEsQuery {
			continue
		}
		if statusQuery != nil && version.Status != *statusQuery {
			continue
		}
		if outdatedQuery != nil && version.Outdated != *outdatedQuery {
			continue
		}
		if vulnerableQuery != nil && version.Vulnerable != *vulnerableQuery {
			continue
		}
		responsePayload = append(responsePayload, Version{
			Name:           version.Name,
			ServiceID:      serviceID,
			CurrentVersion: version.CurrentVersion,
			CurrentCVEs:    version.CurrentCVEs,
			LatestVersion:  version.LatestVersion,
			LatestCVEs:     version.LatestCVEs,
			Status:         schema.ServiceVersionStatusNames[version.Status],
			Outdated:       version.Outdated,
			Vulnerable:     version.Vulnerable,
		})
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func GetServicesByIDVersionsByNameDetails(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	serviceID := req.PathValue("serviceID")
	versionName := req.PathValue("versionName")
	if serviceID == "" || versionName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	details := app.GetServiceVersionDetails(serviceID, versionName)
	if details == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	responsePayload := VersionDetails{
		CurrentVersionNotes: details.CurrentVersionNotes,
		LatestVersionNotes:  details.LatestVersionNotes,
	}
	responsePayload.CurrentCVEDescriptions = make([]CVE, len(details.CurrentCVEDescriptions))
	for i, cveDescription := range details.CurrentCVEDescriptions {
		responsePayload.CurrentCVEDescriptions[i] = CVE{
			Name:        cveDescription.Name,
			Description: cveDescription.Description,
			URL:         cveDescription.URL,
		}
	}
	responsePayload.LatestCVEDescriptions = make([]CVE, len(details.LatestCVEDescriptions))
	for i, cveDescription := range details.LatestCVEDescriptions {
		responsePayload.LatestCVEDescriptions[i] = CVE{
			Name:        cveDescription.Name,
			Description: cveDescription.Description,
			URL:         cveDescription.URL,
		}
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}
