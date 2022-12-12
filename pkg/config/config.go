package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type Configuration struct {
	Profile     string
	Application ApplicationConfiguration
	Http        HttpConfiguration
	Database    map[string](DatabaseConnectionConfiguration)
	Logging     LogConfiguration
	Amazonses   AmazonSESConfiguration
}

type ApplicationConfiguration struct {
	Name      string
	JWTSecret string
}

type AmazonSESConfiguration struct {
	PasswordResetUrl string
	Sender           string
	AccessKeyID      string
	SecretAccessKey  string
}

type HttpConfiguration struct {
	PortNumber     int
	ReadTimeout    int
	WriteTimeout   int
	MaxHeaderBytes int
}

type LogConfiguration struct {
	LogFile map[string](LogFileConfiguration)
}

type LogFileConfiguration struct {
	LogLevel string
	Path     string
	Name     string
}

type DatabaseConnectionConfiguration struct {
	Driver       string
	HostName     string
	PortNumber   uint16
	DatabaseName string
	UserName     string
	Password     string
}

var ServerConfiguration *Configuration

// init loads the configuration YAML file to initialize the
// key/value pairs.
func init() {
	env := os.Getenv("APP_ENV_PROFILE")
	if len(env) == 0 {
		env = "dev"
	}

	configFilePath := filepath.Clean(filepath.Join(".", "conf", "config-"+env+".yml"))
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	var currentConfig Configuration
	if err := yaml.Unmarshal(file, &currentConfig); err != nil {
		panic(err)
	}

	ServerConfiguration = &currentConfig
}
