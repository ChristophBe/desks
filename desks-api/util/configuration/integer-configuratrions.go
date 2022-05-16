package configuration

type IntegerConfigurationKey string

func (key IntegerConfigurationKey) GetValue() int {
	return integerConfigurations[key]
}

var integerConfigurations = make(map[IntegerConfigurationKey]int)

func RegisterIntegerConfigurationWithDefault(key IntegerConfigurationKey, defaultValue int) {

	integerConfigurations[key] = getIntEnvironmentVariable(string(key), defaultValue)
}
