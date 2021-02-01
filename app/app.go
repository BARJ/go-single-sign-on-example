package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Start() {
	config := struct {
		APIPort              int    `env:"API_PORT" default:"8080"`
		DBHost               string `env:"DB_HOST" default:"sso-database"`
		DBPort               int    `env:"DB_PORT" default:"5432"`
		DBUser               string `env:"DB_USER" default:"sso"`
		DBPassword           string `env:"DB_PASSWORD" default:"sso"`
		DBName               string `env:"DB_NAME" default:"sso"`
		JWTSecret            string `env:"JWT_SECRET"`
		AuthTokenLifetime    string `env:"AUTH_TOKEN_LIFETIME" default:"604800s"` // 1 week
		GoogleClientID       string `env:"GOOGLE_CLIENT_ID" default:""`
		GoogleClientSecret   string `env:"GOOGLE_CLIENT_SECRET" default:""`
		GoogleRedirectURI    string `env:"GOOGLE_REDIRECT_URI" default:"https://localhost/api/v1/single-sign-on/google/callback"`
		FacebookClientID     string `env:"FACEBOOK_CLIENT_ID" default:""`
		FacebookClientSecret string `env:"FACEBOOK_CLIENT_SECRET" default:""`
		FacebookRedirectURI  string `env:"FACEBOOK_REDIRECT_URI" default:"https://localhost/api/v1/single-sign-on/facebook/callback"`
		GithubClientID       string `env:"GITHUB_CLIENT_ID" default:""`
		GithubClientSecret   string `env:"GITHUB_CLIENT_SECRET" default:""`
		GithubRedirectURI    string `env:"GITHUB_REDIRECT_URI" default:"https://localhost/api/v1/single-sign-on/github/callback"`
		HomepageURL          string `env:"HOMEPAGE_URL" default:"https://localhost"`
	}{}
	err := NewEnv().Load(&config)
	if err != nil {
		log.Fatal(err)
	}

	database := NewDatabase(DBConfig{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
	})
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.MigrateUp(db); err != nil {
		log.Fatal(err)
	}

	repository := NewSqlRepository(db)
	tokenizer := NewJWT([]byte(config.JWTSecret))
	authenticator := NewAuthenticator(tokenizer, time.Duration(7*24*time.Hour))
	singleSignOnFactory := NewSingleSignOnFactory(authenticator, repository)
	googleSingleSignOn := singleSignOnFactory.NewSingleSignOn(
		NewGoogleIdentityProvider(
			config.GoogleClientID,
			config.GoogleClientSecret,
			config.GoogleRedirectURI,
		),
	)
	facebookSingleSignOn := singleSignOnFactory.NewSingleSignOn(
		NewFacebookIdentityProvider(
			config.FacebookClientID,
			config.FacebookClientSecret,
			config.FacebookRedirectURI,
		),
	)
	githubSingleSignOn := singleSignOnFactory.NewSingleSignOn(
		NewGithubIdentityProvider(
			config.GithubClientID,
			config.GithubClientSecret,
			config.GithubRedirectURI,
		),
	)

	router := NewRouter(
		config.HomepageURL,
		NewStatusHandler(),
		NewUserHandler(authenticator, NewUserManager(repository)),
		googleSingleSignOn,
		facebookSingleSignOn,
		githubSingleSignOn,
	)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), router))
}
