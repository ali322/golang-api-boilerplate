package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var App = new(AppConf)
var Database = new(DatabaseConf)

type AppConf struct {
	Port           string `yaml:"port"`
	JWTSecret      string `yaml:"jwtSecret"`
	Locale         string `yaml:"locale"`
	LogDir         string `yaml:"logDir"`
	GroupAdminRole string `yaml:"groupAdminRole"`
	DefaultRole    string `yaml:"defaultRole"`
	DatabaseURL    string `yaml:"url"`
}

type DatabaseConf struct {
	URL string `yaml:"url"`
}

func Read() {
	workDir, _ := os.Getwd()
	viper.SetConfigFile(filepath.Join(workDir, "config.yml"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Sub("app").Unmarshal(App); err != nil {
		log.Fatal(err)
	}
	if err := viper.Sub("database").Unmarshal(Database); err != nil {
		log.Fatal(err)
	}
}
