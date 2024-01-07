package database

import (
	"fmt"
	"lt/app/models"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBhelper struct {
	log *zap.SugaredLogger
}

func NewDB(log *zap.SugaredLogger) *DBhelper {
	return &DBhelper{
		log: log,
	}
}

func (d *DBhelper) InitDB() *gorm.DB {
	dbname := viper.GetString("DB_NAME")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dialector := postgres.New(postgres.Config{
		DSN:                  dbInfo,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.InternalTrack{})

	return db
}
