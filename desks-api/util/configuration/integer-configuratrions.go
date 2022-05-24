package configuration

import "log"

type IntegerConfigurationKey string

func (key IntegerConfigurationKey) GetValue() int {
	if configValue, ok := integerConfigurations[key]; ok {
		return configValue
	}
	log.Fatalf("%s not found in configurations", key)
	return 0
}

var integerConfigurations = make(map[IntegerConfigurationKey]int)

func RegisterIntegerConfigurationWithDefault(key IntegerConfigurationKey, defaultValue int) {

	integerConfigurations[key] = getIntEnvironmentVariable(string(key), defaultValue)
}
