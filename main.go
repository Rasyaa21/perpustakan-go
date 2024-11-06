package main

import (
	"base-gin/config"
	_ "base-gin/docs"
	"base-gin/repository"
	"base-gin/rest"
	"base-gin/server"
	"base-gin/service"
	"base-gin/storage"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Base API Service
//	@version		1.0
//	@description	This is a boilerplate project, please update accordingly.

//	@contact.name	Mark Muhammad
//	@contact.email	mark.p.e.muhammad@gmail.com

//	@license.name	MIT

//	@host		localhost:3000
//	@BasePath	/v1

//	@securityDefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer auth containing JWT

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg := config.NewConfig()
	storage.InitDB(cfg)
	repository.SetupRepositories()
	service.SetupServices(&cfg)

	app := server.Init(&cfg, repository.GetAccountRepo())
	rest.SetupRestHandlers(app)

	// Swagger
	if cfg.App.Mode == "debug" {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	server.Serve(app)
}
