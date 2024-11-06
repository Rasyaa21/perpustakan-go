package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type AppConfig struct {
	Name    string `env:"APP_NAME"`
	Address string `env:"SERVER_ADDRESS"`
	Mode    string `env:"GIN_MODE" envDefault:"release"`
}

type DBConfig struct {
	DSN           string `env:"DB_DSN"`
	MaxOpenPool   int    `env:"DB_MAX_OPEN_POOL" envDefault:"25"`
	MaxIdlePool   int    `env:"DB_MAX_IDLE_POOL" envDefault:"25"`
	MaxIdleSecond int    `env:"DB_MAX_IDLE_SECOND" envDefault:"300"`
}

type AuthNConfig struct {
	LoginThrottleTTL         int    `env:"LOGIN_THROTTLE_TTL" envDefault:"300"` // in seconds
	LoginMaxAttempt          int    `env:"LOGIN_MAX_ATTEMPT" envDefault:"10"`
	JWTSecretKey             string `env:"JWT_SECRET"`
	JWTAuthTTL               int    `env:"JWT_AUTH_TTL" envDefault:"3600"`
	JWTRefreshTTL            int    `env:"JWT_REFRESH_TTL" envDefault:"2592000"`
	PasswordEncryptionSecret string `env:"PWD_SECRET_32CHAR"`
}

type Config struct {
	App   AppConfig
	DB    DBConfig
	AuthN AuthNConfig
}

func NewConfig() Config {
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if os.Getenv("SERVER_ADDRESS") == "" {
		log.Info().Msg("OS Env not found. Load .env file")
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("Failed to load .env")
		}
	} else {
		log.Info().Msg("OS Env found")
	}

	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	if len(cfg.AuthN.PasswordEncryptionSecret) != 32 {
		log.Fatal().Err(fmt.Errorf("PWD_SECRET_32CHAR must be %d characters", 32)).Msg("config error")
	}

	return cfg
}
