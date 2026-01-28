package api

import (
	"encoding/json"
	"net/http"
)

func GetUser(res http.ResponseWriter, req *http.Request) {
	authorizedUser := GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	json.NewEncoder(res).Encode(UserFromData(*authorizedUser))
}
