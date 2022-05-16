package configuration

type StringConfigurationKey string

func (key StringConfigurationKey) GetValue() string {
	return stringConfigurations[key]
}

var stringConfigurations = make(map[StringConfigurationKey]string)

func RegisterStringConfigurationWithDefault(key StringConfigurationKey, defaultValue string) {

	stringConfigurations[key] = getStringEnvironmentVariable(string(key), defaultValue)
}
