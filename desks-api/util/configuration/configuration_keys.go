package configuration

const (
	ServerPort = IntegerConfigurationKey("SERVER_PORT")
	BaseUrl    = StringConfigurationKey("BASE_URL")

	DBHost     = StringConfigurationKey("DB_HOST")
	DBUsername = StringConfigurationKey("DB_USERNAME")
	DBPassword = StringConfigurationKey("DB_PASSWORD")
	DBName     = StringConfigurationKey("DB_NAME")

	JwtSigningKey                = StringConfigurationKey("JWT_SIGNING_KEY")
	UserdataCleanerIntervalHours = IntegerConfigurationKey("USERDATA_CLEANER_INTERVAL_HOURS")
	MaxUserdataAgeDays           = IntegerConfigurationKey("MAX_USERDATA_AGE_DAYS")

	OauthClientId     = StringConfigurationKey("OAUTH_CLIENT_ID")
	OauthClientSecret = StringConfigurationKey("OAUTH_CLIENT_SECRET")
	OICDIssuerUrl     = StringConfigurationKey("OICD_ISSUER_URL")
)
