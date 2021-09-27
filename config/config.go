package config

import (
	"fmt"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Application = new(appConfig)

type LogConfig struct {
	Prefix   string `json:"prefix"`
	Level    string `json:"level"`
	SaveFile bool   `json:"saveFile"`
}

type JWTConfig struct {
	SigningKey      string `json:"signingKey" yaml:"signingKey"`
	TokenExpiration int64  `json:"tokenExpiration",yaml:"tokenExpiration"`
}

type MongoDBConfig struct {
	ConnURI           string `json:"connURI"`
	DBName            string `json:"dbName"`
	UsersCollection   string `json:"usersCollection"`
	PluginsCollection string `json:"pluginsCollection"`
}

type appConfig struct {
	Port    string        `json:"port"`
	RunMode string        `json:"runMode"`
	Log     LogConfig     `json:"log"`
	Jwt     JWTConfig     `json:"jwt"`
	Mongo   MongoDBConfig `json:"mongo"`
}

func (s appConfig) validate() error {
	return nil
}

func InitConfig() error {
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(path.Join(workDir, "/config"))
	if err := vp.ReadInConfig(); err != nil {
		return err
	}
	vp.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		if err := vp.Unmarshal(Application); err != nil {
			fmt.Println(err)
		}
	})

	if err := vp.Unmarshal(Application); err != nil {
		return err
	}

	return Application.validate()
}
