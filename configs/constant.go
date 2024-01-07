package configs

import "golang.org/x/oauth2/spotify"

var (
	HTTP_HOST string
	HTTP_PORT int
	LOG_LEVEL string

	REDIS_HOST           string
	REDIS_PORT           int
	ENABLE_REDIS_CLUSTER bool

	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_NAME     string
	DB_PASSWORD string

	ACCESS_TOKEN string

	CLIENT_ID     = "3e5d8d3a7e75498aab26fbca864b13e6"
	CLIENT_SECRET = "10ea78148b164aaa9201113874cd351d"
	REDIRECT_URI  = "http://localhost:8888/api/v1/lt/spotify/callback"
	SCOPES        = []string{"user-read-private", "user-read-email"}
	ENDPOINT      = spotify.Endpoint
	STATE         string
)
