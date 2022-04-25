package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func Setup(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mock-server-go" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".mock-server"))
		viper.SetConfigName(".mock-server-go")
	}

	SetDefaults()
	MapEnvVars()

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

}

func SetDefaults() {
	viper.SetDefault("mock.db.engine", "sqlite3")
	viper.SetDefault("mock.db.url", "./mock.db")
}

func MapEnvVars() {
	viper.SetEnvPrefix("mockapp")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
