package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/server/scripts"
)

type Actions struct {
	Groups  []ActionGroup `json:"groups"`
	Actions []Action      `json:"actions"`
}

type ActionGroup struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Action struct {
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	URL              string `json:"url"`
	CanRun           bool   `json:"canRun"`
	Group            string `json:"group"`
	RestrictToAdmins bool   `json:"restrictToAdmins"`
}

func GetServicesByIDActions(res http.ResponseWriter, req *http.Request) {
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

	actionGroups := app.ListServiceActionGroups(serviceID)
	actions := app.ListServiceActions(serviceID)

	var responsePayload Actions
	responsePayload.Groups = make([]ActionGroup, len(actionGroups))
	for i, actionGroup := range actionGroups {
		responsePayload.Groups[i].Name = actionGroup.Name
		responsePayload.Groups[i].Icon = actionGroup.Icon
	}
	responsePayload.Actions = make([]Action, 0, len(actions))
	for _, action := range actions {
		if action.RestrictToAdmins && authorizedUser.Role != schema.UserRoleAdmin {
			continue
		}
		responsePayload.Actions = append(responsePayload.Actions, Action{
			Name:             action.Name,
			Icon:             action.Icon,
			URL:              action.URL,
			CanRun:           action.Script != "",
			Group:            action.Group,
			RestrictToAdmins: action.RestrictToAdmins,
		})
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(responsePayload)
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func PostServicesByIDActionsByNameRun(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil || authorizedUser.Name == "" {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	serviceID := req.PathValue("serviceID")
	actionName := req.PathValue("actionName")
	if serviceID == "" || actionName == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	action := app.GetServiceAction(serviceID, actionName)
	if action == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if action.RestrictToAdmins {
		if authorizedUser.Role != schema.UserRoleAdmin {
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
