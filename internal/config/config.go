package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string `mapstructure:"PORT"`
	MongoDBURL     string `mapstructure:"MONGO_URI"`
	GinMode        string `mapstructure:"GIN_MODE"`
	DBName         string `mapstructure:"DB_NAME"`
	CollectionName string `mapstructure:"COLLECTION_NAME"`
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
