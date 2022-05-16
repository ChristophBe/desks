package configuration

import (
	"log"
	"os"
	"reflect"
	"strconv"
)

func getStringEnvironmentVariable(key string, defaultValue string) string {
	value := os.Getenv(key)
	if reflect.ValueOf(value).IsZero() {
		return defaultValue
	}
	return value
}

func getIntEnvironmentVariable(key string, defaultValue int) int {
	valueString := os.Getenv(key)

	if reflect.ValueOf(valueString).IsZero() {
		return defaultValue
	}

	value, err := strconv.Atoi(valueString)
	if err != nil {
		log.Fatalln("Can not convert environment variable " + key + " to string.")
		return defaultValue
	}

	if reflect.ValueOf(value).IsZero() {
		return defaultValue
	}
	return value
}
