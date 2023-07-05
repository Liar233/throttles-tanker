package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/Liar233/throttles-tank/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger = logrus.New()

func main() {

	if len(os.Args) != 2 {

		Logger.Errorln("config file not set")
		os.Exit(1)
	}

	p, err := filepath.Abs(os.Args[1])

	if err != nil {

		Logger.Errorf("config file path not valid: %s\n", err.Error())
		os.Exit(1)
	}

	dir, filename := path.Split(p)

	conf, err := loadConfig(dir, filename)

	if err != nil {

		Logger.Error("failed parsing config file: %s\n", err.Error())
		os.Exit(1)
	}

	appServer := server.NewAppServer(*conf)

	appServer.Bootstrap()

	if err = appServer.Run(); err != nil {

		Logger.Errorf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func loadConfig(path, filename string) (*server.AppServerConfig, error) {

	viper.SetConfigType("yaml")
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {

		return nil, errors.New("failed read config")
	}

	appConfig := &server.AppServerConfig{}

	if err := viper.Unmarshal(appConfig); err != nil {

		return nil, err
	}

	return appConfig, nil
}
