package app

import (
	"lt/api"
	"lt/app/client/cache"
	"lt/app/client/database"
	"lt/app/lib"
	"lt/app/modules/lt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var LOGGER *zap.SugaredLogger
var ControllerRegistry = map[string]api.HTTPController{}

func Start(cmd *cobra.Command, args []string) {
	// Making directory for logs
	os.Mkdir("logs", os.ModePerm)

	LOGGER = lib.SetUpLogger()
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetInt("REDIS_PORT")

	redisCache := cache.NewRedisCache(LOGGER, redisHost, redisPort)

	err := redisCache.Connect()
	if err != nil {
		LOGGER.Info(err)
		os.Exit(-1)
	}

	LOGGER.Info("Redis started.")

	dbClient := database.NewDB(LOGGER)
	dbInstance := dbClient.InitDB()

	ltMgr := lt.NewManager(LOGGER, dbInstance, redisCache)

	ControllerRegistry["lt"] = lt.NewController(LOGGER, ltMgr)

	LOGGER.Info("Spotify Service  started.")

	initHTTP()
	ch := make(chan os.Signal)
	<-ch
}

func initHTTP() {
	httpHost := viper.GetString("HTTP_HOST")
	httpPort := viper.GetString("HTTP_PORT")
	webServer := api.NewWebServer(LOGGER, httpHost, httpPort)
	for _, value := range ControllerRegistry {
		webServer.SetRoute(value)
	}

	go func() {
		err := webServer.Start()
		if err != nil {
			LOGGER.Infof("Error in web server %v", err)
		}
		LOGGER.Infof("Web server closed")
	}()
}
