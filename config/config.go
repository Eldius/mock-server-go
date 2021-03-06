package config

import (
	"github.com/spf13/viper"
)

func GetDbEngine() string {
	return viper.GetString("mock.db.engine")
}

func GetDbUrl() string {
	return viper.GetString("mock.db.url")
}
