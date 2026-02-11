package state

import (
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/data/schema"
	"github.com/containrrr/shoutrrr"
)

func alertScriptError(app *application.App, script schema.ServiceScript, _ error) {
	adminRole := schema.UserRoleAdmin
	t := true
	f := false
	service, err := app.GetService(schema.ServiceQuery{ID: &script.ServiceID})
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
	}
	channels, err := app.ListServiceUserAlertChannels(
		script.ServiceID,
		&schema.ServiceUserAlertChannelQuery{
			UserQuery:          schema.UserQuery{Role: &adminRole},
			ServiceConfigQuery: schema.ServiceConfigQuery{Hidden: &f},
			ErrorAlerts:        &t,
		},
	)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		return
	}
	channelNames := make(map[string]struct{}, len(channels))
	for _, channel := range channels {
		if channel.URL != "" {
			url, err := crypto.Decrypt(app.Config.SecretKey, channel.URL)
			if err != nil {
				log.Printf("Error decrypting value - %s\n", err)
				continue
			}
			channelNames[url] = struct{}{}
		}
	}
	sender, err := shoutrrr.CreateSender(slices.Collect(maps.Keys(channelNames))...)
	if err != nil {
		log.Printf("Error sending alerts - %s", err)
		return
	}
	var message string
	if service.Name != "" {
		message = fmt.Sprintf("Script \"%s\" for service \"%s\" failed to run", script.Name, service.Name)
	} else {
		message = fmt.Sprintf("Script \"%s\" failed to run", script.Name)
	}
	for _, err := range sender.Send(message, nil) {
		if err != nil {
			log.Printf("Error sending alerts - %s\n", err)
		}
	}
}

func alertUptimeStatuses(app *application.App, script schema.ServiceScript, _ []ServiceUptimeStatusState) {
	t := true
	f := false
	service, err := app.GetService(schema.ServiceQuery{ID: &script.ServiceID})
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
	}
	channels, err := app.ListServiceUserAlertChannels(
		script.ServiceID,
		&schema.ServiceUserAlertChannelQuery{
			ServiceConfigQuery: schema.ServiceConfigQuery{Hidden: &f},
			UptimeAlerts:       &t,
		},
	)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		return
	}
	channelNames := make(map[string]struct{}, len(channels))
	for _, channel := range channels {
		if channel.URL != "" {
			url, err := crypto.Decrypt(app.Config.SecretKey, channel.URL)
			if err != nil {
				log.Printf("Error decrypting value - %s\n", err)
				continue
			}
			channelNames[url] = struct{}{}
		}
	}
	sender, err := shoutrrr.CreateSender(slices.Collect(maps.Keys(channelNames))...)
	if err != nil {
		log.Printf("Error sending alerts - %s", err)
		return
	}
	var message string
	if service.Name != "" {
		message = fmt.Sprintf("Service \"%s\" is down", service.Name)
	} else {
		message = "A service is down"
	}
	for _, err := range sender.Send(message, nil) {
		if err != nil {
			log.Printf("Error sending alerts - %s\n", err)
		}
	}
}

func alertVersions(app *application.App, script schema.ServiceScript, _ []ServiceVersionState) {
	t := true
	f := false
	service, err := app.GetService(schema.ServiceQuery{ID: &script.ServiceID})
	if err != nil {
		log.Printf("Error loading record - %s\n", err)
	}
	channels, err := app.ListServiceUserAlertChannels(
		script.ServiceID,
		&schema.ServiceUserAlertChannelQuery{
			ServiceConfigQuery: schema.ServiceConfigQuery{Hidden: &f},
			VersionAlerts:      &t,
		},
	)
	if err != nil {
		log.Printf("Error loading records - %s\n", err)
		return
	}
	channelNames := make(map[string]struct{}, len(channels))
	for _, channel := range channels {
		if channel.URL != "" {
			url, err := crypto.Decrypt(app.Config.SecretKey, channel.URL)
			if err != nil {
				log.Printf("Error decrypting value - %s\n", err)
				continue
			}
			channelNames[url] = struct{}{}
		}
	}
	sender, err := shoutrrr.CreateSender(slices.Collect(maps.Keys(channelNames))...)
	if err != nil {
		log.Printf("Error sending alerts - %s", err)
		return
	}
	var message string
	if service.Name != "" {
		message = fmt.Sprintf("Service \"%s\" is outdated", service.Name)
	} else {
		message = "A service is outdated"
	}
	for _, err := range sender.Send(message, nil) {
		if err != nil {
			log.Printf("Error sending alerts - %s\n", err)
		}
	}
}
