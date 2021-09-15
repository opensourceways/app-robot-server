package config

import (
	"fmt"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Application = new(server)

type Log struct {
	Prefix   string `json:"prefix" yaml:"prefix"`
	Level    string `json:"level" yaml:"level"`
	SaveFile bool   `json:"saveFile" yaml:"saveFile"`
}

type JWT struct {
	SigningKey      string `json:"signingKey" yaml:"signingKey"`
	TokenExpiration int64  `json:"tokenExpiration",yaml:"tokenExpiration"`
}

type server struct {
	Port    string `json:"port" yaml:"port"`
	RunMode string `json:"runMode" yaml:"runMode"`
	Log     Log    `json:"log" yaml:"log"`
	Jwt     JWT    `json:"jwt" yaml:"jwt"`
}

func (s server) validate() error {
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
