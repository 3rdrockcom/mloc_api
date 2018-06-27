package log

import (
	"github.com/epointpayment/mloc_api_go/app/config"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	SetMode(config.EnvDevelopment)
}

func Start() {
	singletonLogger = &MyLogger{
		Logger: logger,
	}
}

func Stop() {}

func SetMode(env string) {
	switch env {
	case config.EnvProduction:
		setProduction()
	case config.EnvDevelopment:
		setDevelopment()
	}
}

func setDevelopment() {
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &logrus.TextFormatter{}
}

func setProduction() {
	logger.SetLevel(logrus.InfoLevel)
	logger.Formatter = &logrus.JSONFormatter{}
}
