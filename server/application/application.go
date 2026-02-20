package application

import (
	"errors"
	"fmt"
	"log"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data"
	"github.com/2manyvcos/paranal/server/schema"
	"github.com/2manyvcos/paranal/utils"
	"github.com/go-co-op/gocron/v2"
)

type App struct {
	Config struct {
		AppURL      string
		AppName     string
		Tagline     string
		ScriptsPath string
		SecretKey   string

		Brand struct {
			AssetsPath     string
			Logo           string
			Favicon        string
			AppleTouchIcon string
			Theme          string
		}

		Server struct {
			Protocol string
			Address  string
			Port     string
			CertFile string
			KeyFile  string
		}

		DB data.Config

		Admin struct {
			Username     string
			PasswordHash string
		}

		Auth struct {
			LogoutRedirectURL string

			RemoteUser struct {
				Enabled            bool
				HeaderName         string
				GroupsHeaderName   string
				AdminGroup         string
				CreateUnknownUsers bool
				Whitelist          []string
			}

			JWT struct {
				Secret string
			}
		}
	}

	schema.DataProvider
	schema.State
	Scheduler gocron.Scheduler
}

func Setup() (app *App, err error) {
	app = new(App)

	app.Config.AppURL = utils.LoadConfigValue("APP_URL", "")
	app.Config.AppName = utils.LoadConfigValue("APP_NAME", "Paranal")
	app.Config.Tagline = utils.LoadConfigValue("TAGLINE", "Advanced Service Dashboard with Health and Version Monitoring")
	app.Config.ScriptsPath = utils.LoadConfigValue("SCRIPTS_PATH", ".paranal/scripts")
	if app.Config.SecretKey, err = utils.LoadRequiredConfigValue("SECRET_KEY"); err != nil {
		app.Close()
		return nil, err
	}

	app.Config.Brand.AssetsPath = utils.LoadConfigValue("BRAND_ASSETS_PATH", "")
	app.Config.Brand.Logo = utils.LoadConfigValue("BRAND_LOGO", "/logo.svg")
	app.Config.Brand.Favicon = utils.LoadConfigValue("BRAND_FAVICON", "/favicon.svg")
	app.Config.Brand.AppleTouchIcon = utils.LoadConfigValue("BRAND_APPLE_TOUCH_ICON", "/apple-touch-icon.png")
	app.Config.Brand.Theme = utils.LoadConfigValue("BRAND_THEME", "/theme.css")

	app.Config.Server.Protocol = utils.LoadConfigValue("SERVER_PROTOCOL", "http")
	app.Config.Server.Address = utils.LoadConfigValue("SERVER_ADDRESS", "0.0.0.0")
	app.Config.Server.Port = utils.LoadConfigValue("SERVER_PORT", "8080")
	app.Config.Server.CertFile = utils.LoadConfigValue("SERVER_CERT_FILE", ".paranal/ssl/cert.pem")
	app.Config.Server.KeyFile = utils.LoadConfigValue("SERVER_KEY_FILE", ".paranal/ssl/key.pem")

	app.Config.DB.Type = utils.LoadConfigValue("DB_TYPE", "sqlite")
	app.Config.DB.Path = utils.LoadConfigValue("DB_PATH", ".paranal/paranal.db")

	app.Config.Admin.Username = utils.LoadConfigValue("ADMIN_USERNAME", "")
	app.Config.Admin.PasswordHash = utils.LoadConfigValue("ADMIN_PASSWORD_HASH", "")

	app.Config.Auth.LogoutRedirectURL = utils.LoadConfigValue("AUTH_LOGOUT_REDIRECT_URL", "")
	app.Config.Auth.RemoteUser.Enabled = utils.LoadConfigBool("AUTH_REMOTE_USER_ENABLED", false)
	app.Config.Auth.RemoteUser.HeaderName = utils.LoadConfigValue("AUTH_REMOTE_USER_HEADER_NAME", "Remote-User")
	app.Config.Auth.RemoteUser.GroupsHeaderName = utils.LoadConfigValue("AUTH_REMOTE_USER_GROUPS_HEADER_NAME", "")
	app.Config.Auth.RemoteUser.AdminGroup = utils.LoadConfigValue("AUTH_REMOTE_USER_ADMIN_GROUP", "")
	app.Config.Auth.RemoteUser.CreateUnknownUsers = utils.LoadConfigBool("AUTH_REMOTE_USER_CREATE_UNKNOWN_USERS", true)
	app.Config.Auth.RemoteUser.Whitelist = utils.LoadConfigList("AUTH_REMOTE_USER_WHITELIST", nil)
	app.Config.Auth.JWT.Secret = crypto.GenerateJWTTokenSecret()

	app.DataProvider, err = data.Load(app.Config.DB)
	if err != nil {
		app.Close()
		return nil, err
	}

	app.Scheduler, err = gocron.NewScheduler()
	if err != nil {
		app.Close()
		return nil, err
	}

	err = app.prepare()
	if err != nil {
		app.Close()
		return nil, err
	}

	return
}

func (app *App) prepare() (err error) {
	if app.Config.Admin.Username != "" {
		log.Printf("Setting up admin user \"%s\"\n", app.Config.Admin.Username)
		existing, err := app.GetUser(schema.UserQuery{Name: &app.Config.Admin.Username})
		if err != nil && !errors.Is(err, schema.ErrNotFound) {
			return fmt.Errorf("Error loading record - %s\n", err)
		}
		newUser := existing
		newUser.Name = app.Config.Admin.Username
		newUser.PasswordHash = app.Config.Admin.PasswordHash
		newUser.Role = schema.UserRoleAdmin
		err = app.CreateOrUpdateUser(newUser)
		if err != nil {
			return err
		}
	}

	log.Println("Cleaning up database")
	app.CleanupDatabase()

	return nil
}

func (app *App) Close() {
	if app.DataProvider != nil {
		app.DataProvider.Close()
	}

	if app.State != nil {
		app.State.Close()
	}

	if app.Scheduler != nil {
		app.Scheduler.Shutdown()
	}
}
