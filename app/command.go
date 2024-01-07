package app

import (
	"fmt"
	"log"

	"lt/configs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initEnvConfig() {
	var err error

	// SQL DB
	if err = viper.BindEnv("DB_USER"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("DB_PASSWORD"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("DB_NAME"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("DB_HOST"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("DB_PORT"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}

	// Redis
	if err = viper.BindEnv("REDIS_HOST"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("REDIS_PORT"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}

	if err = viper.BindEnv("LOG_LEVEL"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}

	// WEB SERVER
	if err = viper.BindEnv("HTTP_HOST"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("HTTP_PORT"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}

	if err = viper.BindEnv("CLIENT_ID"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("CLIENT_SECRET"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv(""); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("SCOPES"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}
	if err = viper.BindEnv("STATE"); err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}

	if err != nil {
		log.Fatal(fmt.Sprintf("unable to read config values %v", err))
	}

}

func initConfig() {

	configs.LOG_LEVEL = viper.GetString("LOG_LEVEL")

	configs.HTTP_HOST = viper.GetString("HTTP_HOST")
	configs.HTTP_PORT = viper.GetInt("HTTP_PORT")

	configs.REDIS_HOST = viper.GetString("REDIS_HOST")
	configs.REDIS_PORT = viper.GetInt("REDIS_PORT")

	configs.DB_USER = viper.GetString("DB_USER")
	configs.DB_PASSWORD = viper.GetString("DB_PASSWORD")
	configs.DB_HOST = viper.GetString("DB_HOST")
	configs.DB_NAME = viper.GetString("DB_NAME")
	configs.DB_PORT = viper.GetInt("DB_PORT")

	configs.CLIENT_ID = viper.GetString("CLIENT_ID")
	configs.CLIENT_SECRET = viper.GetString("CLIENT_SECRET")
	configs.REDIRECT_URI = viper.GetString("REDIRECT_URI")
	configs.SCOPES = viper.GetStringSlice("SCOPES")

	configs.STATE = viper.GetString("STATE")

}

func init() {
	initEnvConfig()
	cobra.OnInitialize(initConfig)
	viper.SetConfigName("lt")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("No config file found %s,"+" Default values from environment variables. \n", err)
	}

	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	_ = viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	err = viper.SafeWriteConfig()
	if err != nil {
		_ = fmt.Errorf("failed to write config %v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "lt",
	Short: "lt",
	Run:   Start,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
