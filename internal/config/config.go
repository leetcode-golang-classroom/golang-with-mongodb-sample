package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	MongoDBURL         string `mapstructure:"MONGO_URI"`
	GinMode            string `mapstructure:"GIN_MODE"`
	DBName             string `mapstructure:"DB_NAME"`
	CollectionName     string `mapstructure:"COLLECTION_NAME"`
	NewRelicKey        string `mapstructure:"NEW_RELIC_KEY"`
	NewRelicLicenseKey string `mapstructure:"NEW_RELIC_LICENSE_KEY"`
	AppName            string `mapstructure:"APP_NAME"`
	Environment        string `mapstructure:"ENV"`
}

var AppConfig *Config

func init() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	FailOnError(v.BindEnv("PORT"), "failed to bind PORT")
	FailOnError(v.BindEnv("MONGO_URI"), "failed to bind MONGO_URI")
	FailOnError(v.BindEnv("GIN_MODE"), "failed to bind GIN_MODE")
	FailOnError(v.BindEnv("DB_NAME"), "failed to bind DB_NAME")
	FailOnError(v.BindEnv("COLLECTION_NAME"), "failed to bind COLLECTION_NAME")
	FailOnError(v.BindEnv("NEW_RELIC_KEY"), "failed to bind NEW_RELIC_KEY")
	FailOnError(v.BindEnv("NEW_RELIC_LICENSE_KEY"), "failed to bind NEW_RELIC_LICENSE_KEY")
	FailOnError(v.BindEnv("APP_NAME"), "failed to bind APP_NAME")
	FailOnError(v.BindEnv("ENV"), "failed to bind ENV")
	err := v.ReadInConfig()
	if err != nil {
		log.Println("Load from environment variable")
	}
	err = v.Unmarshal(&AppConfig)
	if err != nil {
		FailOnError(err, "Failed to read enivronment")
	}
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
