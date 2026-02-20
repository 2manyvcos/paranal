package api

import (
	"fmt"
	"net/http"

	"github.com/2manyvcos/paranal/server/helper"
)

func GetEvents(res http.ResponseWriter, req *http.Request) {
	app := helper.GetApp(req)
	authorizedUser := helper.GetAuthorizedUser(req)

	if authorizedUser == nil {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(res)

	for event := range app.SubscribeToClientEvents(req.Context()) {
		if user := event.User(); user != "" && authorizedUser.Name != user {
			continue
		}

		if eventType := event.Event(); eventType == "" {
			_, err := fmt.Fprintf(res, "data: %s\n\n", event.Data())
			if err != nil {
				return
			}
		} else {
			_, err := fmt.Fprintf(res, "event: %s\ndata: %s\n\n", eventType, event.Data())
			if err != nil {
				return
			}
		}
		err := rc.Flush()
		if err != nil {
			return
		}
	}
}
