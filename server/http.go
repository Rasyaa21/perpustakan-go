package server

import (
	"base-gin/config"
	"base-gin/repository"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const (
	ParamTokenUser     = "x-token-user"
	ParamTokenUserID   = "x-token-user-id"
	ParamTokenUsername = "x-token-uname"
)

var (
	handler *Handler
)

func Init(
	cfg *config.Config,
	accountRepo *repository.AccountRepository,
) *gin.Engine {
	app := gin.New()
	app.Use(gin.Recovery())       // panic handling
	registerCustomValidationTag() // returns json field name on errors

	handler = NewHandler(cfg, accountRepo)

	return app
}

func registerCustomValidationTag() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func Serve(handler http.Handler) {
	srv := &http.Server{
		Addr:              os.Getenv("SERVER_ADDRESS"),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      100 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Stack().Err(err).Msg("Graceful Errors: Http error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Graceful Info: Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Stack().Err(err).Msg("Graceful Errors: Server forced to shutdown")
	}

	log.Info().Msg("Graceful Info: Server exiting")
}

func GetHandler() *Handler {
	if handler == nil {
		panic("handler is no initialised")
	}

	return handler
}
