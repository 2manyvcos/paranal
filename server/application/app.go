package application

import (
	"log"

	"github.com/2manyvcos/paranal/crypto"
	"github.com/2manyvcos/paranal/server/data"
)

type App struct {
	Config struct {
		AppName string
		Tagline string

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

	data.DataProvider
}

func Setup() (app *App, err error) {
	app = new(App)

	app.Config.AppName = loadConfigValue("APP_NAME", "Paranal")
	app.Config.Tagline = loadConfigValue("TAGLINE", "Advanced Service Dashboard with Health and Version Monitoring")

	app.Config.Brand.AssetsPath = loadConfigValue("BRAND_ASSETS_PATH", "")
	app.Config.Brand.Logo = loadConfigValue("BRAND_LOGO", "/logo.svg")
	app.Config.Brand.Favicon = loadConfigValue("BRAND_FAVICON", "/favicon.svg")
	app.Config.Brand.AppleTouchIcon = loadConfigValue("BRAND_APPLE_TOUCH_ICON", "/apple-touch-icon.png")
	app.Config.Brand.Theme = loadConfigValue("BRAND_THEME", "/theme.css")

	app.Config.Server.Protocol = loadConfigValue("SERVER_PROTOCOL", "http")
	app.Config.Server.Address = loadConfigValue("SERVER_ADDRESS", "0.0.0.0")
	app.Config.Server.Port = loadConfigValue("SERVER_PORT", "8080")
	app.Config.Server.CertFile = loadConfigValue("SERVER_CERT_FILE", ".paranal/ssl/cert.pem")
	app.Config.Server.KeyFile = loadConfigValue("SERVER_KEY_FILE", ".paranal/ssl/key.pem")

	app.Config.DB.Type = loadConfigValue("DB_TYPE", "sqlite")
	app.Config.DB.Path = loadConfigValue("DB_PATH", ".paranal/paranal.db")

	app.Config.Admin.Username = loadConfigValue("ADMIN_USERNAME", "")
	app.Config.Admin.PasswordHash = loadConfigValue("ADMIN_PASSWORD_HASH", "")

	app.Config.Auth.LogoutRedirectURL = loadConfigValue("AUTH_LOGOUT_REDIRECT_URL", "")
	app.Config.Auth.RemoteUser.Enabled = loadConfigBool("AUTH_REMOTE_USER_ENABLED", false)
	app.Config.Auth.RemoteUser.HeaderName = loadConfigValue("AUTH_REMOTE_USER_HEADER_NAME", "Remote-User")
	app.Config.Auth.RemoteUser.GroupsHeaderName = loadConfigValue("AUTH_REMOTE_USER_GROUPS_HEADER_NAME", "")
	app.Config.Auth.RemoteUser.AdminGroup = loadConfigValue("AUTH_REMOTE_USER_ADMIN_GROUP", "")
	app.Config.Auth.RemoteUser.CreateUnknownUsers = loadConfigBool("AUTH_REMOTE_USER_CREATE_UNKNOWN_USERS", true)
	app.Config.Auth.RemoteUser.Whitelist = loadConfigList("AUTH_REMOTE_USER_WHITELIST", nil)
	app.Config.Auth.JWT.Secret = crypto.JWTGenerateTokenSecret()

	app.DataProvider, err = data.Load(app.Config.DB)
	if err != nil {
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
		err = app.CreateUser(
			data.User{
				Name:         app.Config.Admin.Username,
				PasswordHash: app.Config.Admin.PasswordHash,
				Role:         data.USER_ROLE_ADMIN,
			},
			true,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
