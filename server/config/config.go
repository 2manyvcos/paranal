package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/2manyvcos/paranal/utils"
)

type Config struct {
	AppName    string
	CORSOrigin string
	PublicAPI  string
	ScriptDir  string

	Server struct {
		Protocol string
		Address  string
		Port     string
		CertFile string
		KeyFile  string
	}

	DB struct {
		File string
	}
}

func Load() (*Config, error) {
	defaultConfigDir, err := utils.ConfigDir("paranal")
	if err != nil {
		return nil, err
	}

	configDir, err := filepath.Abs(utils.LoadConfigValue("PARANAL_CONFIG_DIR", defaultConfigDir))
	if err != nil {
		return nil, fmt.Errorf("resolving config directory: %v", err)
	}

	var config Config

	config.AppName = utils.LoadConfigValue("PARANAL_APP_NAME", "Paranal")
	config.CORSOrigin = utils.LoadConfigValue("PARANAL_CORS_ORIGIN", "")
	config.PublicAPI = strings.TrimSuffix(utils.LoadConfigValue("PARANAL_PUBLIC_API", "/api"), "/")
	config.ScriptDir = utils.LoadConfigValue("PARANAL_SCRIPT_DIR", filepath.Join(configDir, "scripts"))

	config.Server.Protocol = utils.LoadConfigValue("PARANAL_SERVER_PROTOCOL", "http")
	config.Server.Address = utils.LoadConfigValue("PARANAL_SERVER_ADDRESS", "0.0.0.0")
	config.Server.Port = utils.LoadConfigValue("PARANAL_SERVER_PORT", "8080")
	config.Server.CertFile = utils.LoadConfigValue("PARANAL_SERVER_CERT_FILE", filepath.Join(configDir, "tls", "cert.pem"))
	config.Server.KeyFile = utils.LoadConfigValue("PARANAL_SERVER_KEY_FILE", filepath.Join(configDir, "tls", "key.pem"))

	config.DB.File = utils.LoadConfigValue("PARANAL_DB_FILE", filepath.Join(configDir, "paranal.db"))

	return &config, nil
}
