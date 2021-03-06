package configuration

import "log"

type StringConfigurationKey string

func (key StringConfigurationKey) GetValue() string {

	if configValue, ok := stringConfigurations[key]; ok {
		return configValue
	}
	log.Fatalf("%s not found in configurations", key)
	return ""
}

var stringConfigurations = make(map[StringConfigurationKey]string)

func RegisterStringConfigurationWithDefault(key StringConfigurationKey, defaultValue string) {
	stringConfigurations[key] = getStringEnvironmentVariable(string(key), defaultValue)
}

func RegisterStringConfiguration(key StringConfigurationKey) {
	var err error
	stringConfigurations[key], err = getRequiredStringEnvironmentVariable(string(key))
	if err != nil {
		log.Fatal(err)
	}
	return
}
