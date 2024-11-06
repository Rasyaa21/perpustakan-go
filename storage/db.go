package storage

import (
	"base-gin/config"
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(config config.Config) {
	logLevel := logger.Silent
	if config.App.Mode == "debug" {
		logLevel = logger.Error
	}

	zeroWriter := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zeroLogger := logger.New(
		&zeroWriter,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	gormDB, err := gorm.Open(mysql.Open(config.DB.DSN), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		SkipDefaultTransaction: true,
		Logger:                 zeroLogger,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("tidak dapat mengubah pengaturan database")
		log.Fatal().Stack().Err(err).Msg("tidak dapat terhubung ke database")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("tidak dapat terhubung ke database")
	}

	sqlDB.SetMaxOpenConns(config.DB.MaxOpenPool)
	sqlDB.SetMaxIdleConns(config.DB.MaxIdlePool)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DB.MaxIdleSecond) * time.Second)

	db = gormDB
}

func NewDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("db is not initialised")
	}

	return db
}
