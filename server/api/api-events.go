package api

import (
	"fmt"
	"net/http"
)

func GetEvents(res http.ResponseWriter, req *http.Request) {
	app := GetApp(req)

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(res)

	for events := range app.EventBroker.Subscribe(req.Context().Done()) {
		for _, event := range events {
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
		}

		err := rc.Flush()
		if err != nil {
			return
		}
	}
}
