package configuration

const (
	ServerPort = IntegerConfigurationKey("SERVER_PORT")

	DBHost     = StringConfigurationKey("DB_HOST")
	DBUsername = StringConfigurationKey("DB_USERNAME")
	DBPassword = StringConfigurationKey("DB_PASSWORD")
	DBName     = StringConfigurationKey("DB_NAME")

	JwtSigningKey = StringConfigurationKey("JWT_SIGNING_KEY")

	UserdataCleanerIntervalHours = IntegerConfigurationKey("USERDATA_CLEANER_INTERVAL_HOURS")
	MaxUserdataAgeDays           = IntegerConfigurationKey("MAX_USERDATA_AGE_DAYS")
)
