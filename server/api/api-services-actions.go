package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/scripts"
)

func GetServicesByIDActions(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	serviceID := req.PathValue("serviceID")
	if serviceID == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	actionGroups := app.State.GetServiceActionGroups(serviceID)
	actions := app.State.GetServiceActions(serviceID)

	var responsePayload struct {
		Groups []struct {
			Name string `json:"name"`
			Icon string `json:"icon"`
		} `json:"groups"`
		Actions []struct {
			Name             string `json:"name"`
			Icon             string `json:"icon"`
			URL              string `json:"url"`
			CanRun           bool   `json:"canRun"`
			Group            string `json:"group"`
			RestrictToAdmins bool   `json:"restrictToAdmins"`
		} `json:"actions"`
	}
	responsePayload.Groups = make([]struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}, len(actionGroups))
	for i, actionGroup := range actionGroups {
		responsePayload.Groups[i].Name = actionGroup.Name
		responsePayload.Groups[i].Icon = actionGroup.Icon
	}
	responsePayload.Actions = make([]struct {
		Name             string `json:"name"`
		Icon             string `json:"icon"`
		URL              string `json:"url"`
		CanRun           bool   `json:"canRun"`
		Group            string `json:"group"`
		RestrictToAdmins bool   `json:"restrictToAdmins"`
	}, len(actions))
	for i, action := range actions {
		responsePayload.Actions[i].Name = action.Name
		responsePayload.Actions[i].Icon = action.Icon
		responsePayload.Actions[i].URL = action.URL
		responsePayload.Actions[i].CanRun = action.Script != ""
		responsePayload.Actions[i].Group = action.Group
		responsePayload.Actions[i].RestrictToAdmins = action.RestrictToAdmins
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PostServicesByIDActionsRunByName(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	serviceID := req.PathValue("serviceID")
	actionName := req.PathValue("actionName")
	if serviceID == "" || actionName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	action := app.State.GetServiceAction(serviceID, actionName)
	if action == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if action.RestrictToAdmins {
		if authorizedUser == nil || authorizedUser.Role != schema.UserRoleAdmin {
			http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}
	if action.Script == "" {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	results, err := scripts.RunServiceScript(app, serviceID, action.Script)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(res).Encode(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{
			Success: false,
			Error:   err.Error(),
		})
		if err != nil {
			log.Printf("Error encoding response payload - %s\n", err)
		}
		return
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		Success bool  `json:"success"`
		Results []any `json:"results"`
	}{
		Success: true,
		Results: results,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}
