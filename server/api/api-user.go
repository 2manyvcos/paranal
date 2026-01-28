package api

import (
	"encoding/json"
	"net/http"

	"github.com/2manyvcos/paranal/server/helper"
)

func GetUser(res http.ResponseWriter, req *http.Request) {
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	json.NewEncoder(res).Encode(UserFromData(*authorizedUser))
}
