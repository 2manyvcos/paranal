package application

import (
	"strings"

	"github.com/2manyvcos/paranal/server/event"
	"github.com/2manyvcos/paranal/utils"
	"github.com/go-co-op/gocron/v2"
)

type App struct {
	Config struct {
		AppName    string
		CORSOrigin string
		PublicAPI  string
		ScriptPath string

		Server struct {
			Protocol string
			Address  string
			Port     string
			CertFile string
			KeyFile  string
		}
	}

	CronScheduler gocron.Scheduler

	Events *utils.Broker[event.Event]
}

func New() (app *App, err error) {
	app = new(App)

	app.Config.AppName = utils.LoadConfigValue("PARANAL_APP_NAME", "Paranal")
	app.Config.CORSOrigin = utils.LoadConfigValue("PARANAL_CORS_ORIGIN", "")
	app.Config.PublicAPI = strings.TrimSuffix(utils.LoadConfigValue("PARANAL_PUBLIC_API", "/api"), "/")
	app.Config.ScriptPath = utils.LoadConfigValue("PARANAL_SCRIPT_PATH", "/etc/paranal/scripts")

	app.Config.Server.Protocol = utils.LoadConfigValue("PARANAL_SERVER_PROTOCOL", "http")
	app.Config.Server.Address = utils.LoadConfigValue("PARANAL_SERVER_ADDRESS", "0.0.0.0")
	app.Config.Server.Port = utils.LoadConfigValue("PARANAL_SERVER_PORT", "8080")
	app.Config.Server.CertFile = utils.LoadConfigValue("PARANAL_SERVER_CERT_FILE", "cert.pem")
	app.Config.Server.KeyFile = utils.LoadConfigValue("PARANAL_SERVER_KEY_FILE", "key.pem")

	app.CronScheduler, err = gocron.NewScheduler()
	if err != nil {
		// app.Shutdown()
		return nil, err
	}

	app.Events = utils.NewBroker[event.Event]()

	return
}

func (app *App) Shutdown() {
	if app.CronScheduler != nil {
		app.CronScheduler.Shutdown()
	}

	if app.Events != nil {
		app.Events.Shutdown()
	}
}
