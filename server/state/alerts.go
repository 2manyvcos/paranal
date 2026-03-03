package state

import (
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/application"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/containrrr/shoutrrr"
)

func alertScriptError(app *application.App, service *schema.Service, script schema.ServiceScript, _ error) {
	adminRole := schema.UserRoleAdmin
	t := true
	f := false
	channels, err := app.ListServiceAlertChannels(
		script.ServiceID,
		&schema.ServiceAlertChannelQuery{
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
	if service != nil {
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

func alertUnhealthyHealthStatuses(app *application.App, service schema.Service, script schema.ServiceScript, _ []ServiceHealthStatus) {
	t := true
	f := false
	channels, err := app.ListServiceAlertChannels(
		script.ServiceID,
		&schema.ServiceAlertChannelQuery{
			ServiceConfigQuery: schema.ServiceConfigQuery{Hidden: &f},
			HealthAlerts:       &t,
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
	for _, err := range sender.Send(fmt.Sprintf("Service \"%s\" is unhealthy", service.Name), nil) {
		if err != nil {
			log.Printf("Error sending alerts - %s\n", err)
		}
	}
}

func alertOutdatedVersions(app *application.App, service schema.Service, script schema.ServiceScript, _ []ServiceVersion) {
	t := true
	f := false
	channels, err := app.ListServiceAlertChannels(
		script.ServiceID,
		&schema.ServiceAlertChannelQuery{
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
	for _, err := range sender.Send(fmt.Sprintf("Service \"%s\" is outdated", service.Name), nil) {
		if err != nil {
			log.Printf("Error sending alerts - %s\n", err)
		}
	}
}

func alertVulnerableVersions(app *application.App, service schema.Service, script schema.ServiceScript, _ []ServiceVersion) {
	t := true
	f := false
	channels, err := app.ListServiceAlertChannels(
		script.ServiceID,
		&schema.ServiceAlertChannelQuery{
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
	for _, err := range sender.Send(fmt.Sprintf("Service \"%s\" has open CVEs", service.Name), nil) {
		if err != nil {
			log.Printf("Error sending alerts - %s\n", err)
		}
	}
}
