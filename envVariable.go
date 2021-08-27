package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func viperReadEnvVariable(key string) string {
	viper.SetConfigName("database")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
