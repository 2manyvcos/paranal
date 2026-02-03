package api

import (
	"encoding/json"
	"net/http"

	"github.com/2manyvcos/paranal/server/helper"
	"github.com/2manyvcos/paranal/server/scripts"
	"github.com/2manyvcos/paranal/utils"
)

func PostScript(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)

	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		Source    string `json:"source"`
		ServiceID string `json:"serviceID"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	results, err := scripts.Run(app, requestPayload.ServiceID, requestPayload.Source)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(struct {
		Success bool  `json:"success"`
		Results []any `json:"results"`
	}{
		Success: true,
		Results: results,
	})
}
