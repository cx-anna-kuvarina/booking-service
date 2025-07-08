package main

import (
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"time"
)

const (
	httpPortEnv   = "HTTP_PORT"
	dbWriteURLEnv = "DATABASE_WRITE_URL"
	dbReadURLEnv  = "DATABASE_READ_URL"

	googleClientIDEnv     = "GOOGLE_CLIENT_ID"
	googleClientSecretEnv = "GOOGLE_CLIENT_SECRET"
	googleRandomStateEnv  = "GOOGLE_RANDOM_STATE"
	jwtSecretEnv          = "JWT_SECRET"
	jwtExpPeriodEnv       = "JWT_EXP_PERIOD_DURATION"
	appURLEnv             = "APP_URL"
)

const (
	appURlDefault       = "http://localhost"
	jwtSecretDefault    = "abrakcskwq1323dns2"
	jwtExpPeriodDefault = time.Hour
)

var (
	googleScopes = []string{"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile"}
)

type Config struct {
	Port              int
	DBWriteURL        string
	DBReadURL         string
	GoogleLoginConfig oauth2.Config
	GoogleRandomState string
	JWTSecret         string
	JWTExpPeriod      time.Duration
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	googleRedirectURLDefault := fmt.Sprintf("%s:%d", appURlDefault, viper.GetInt(httpPortEnv))
	viper.SetDefault(appURLEnv, googleRedirectURLDefault)
	googleRedirectURL := fmt.Sprintf("%s/google_callback", viper.GetString(appURLEnv))

	viper.SetDefault(jwtSecretEnv, jwtSecretDefault)
	viper.SetDefault(jwtExpPeriodEnv, jwtExpPeriodDefault)

	return &Config{
		Port:       viper.GetInt(httpPortEnv),
		DBWriteURL: viper.GetString(dbWriteURLEnv),
		DBReadURL:  viper.GetString(dbReadURLEnv),
		GoogleLoginConfig: oauth2.Config{
			ClientID:     viper.GetString(googleClientIDEnv),
			ClientSecret: viper.GetString(googleClientSecretEnv),
			Endpoint:     google.Endpoint,
			RedirectURL:  googleRedirectURL,
			Scopes:       googleScopes,
		},
		GoogleRandomState: viper.GetString(googleRandomStateEnv),
		JWTSecret:         viper.GetString(jwtSecretEnv),
		JWTExpPeriod:      viper.GetDuration(jwtExpPeriodEnv),
	}
}
