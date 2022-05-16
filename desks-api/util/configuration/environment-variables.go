package util

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

func getStringEnvironmentVariable(key string, defaultValue string) string {
	value, err := getRequiredStringEnvironmentVariable(key)
	if err != nil {
		return defaultValue
	}
	return value
}
func getRequiredStringEnvironmentVariable(key string) (string, error) {
	value := os.Getenv(key)
	if reflect.ValueOf(value).IsZero() {
		return "", fmt.Errorf("reqwuired environment valriable %s not set", key)
	}
	return value, nil
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
