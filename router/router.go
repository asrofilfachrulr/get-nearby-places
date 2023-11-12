package router

import (
	"net/http"

	"github.com/asrofilfachrulr/get-nearby-places/controller"
	"github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const SWAGGER_URL_REMOTE = "https://gist.githubusercontent.com/asrofilfachrulr/1c8d2a3bf151678c8e01071facb553ed/raw/bbae3a3a2cdc5f6e2f76565599b404dde75db693/get-nearby-places.yaml"

// const SWAGGER_LOCAL = "./swagger-config.yaml"

func SetupRouter(data models.BatchPlace) *gin.Engine {
	r := gin.Default()

	// attach BatchPlace to gin Context
	r.Use(func(ctx *gin.Context) {
		ctx.Set("places", data)
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})

	r.GET("/search", controller.GetNearby)

	// load swagger from remote url
	conf := ginSwagger.URL(SWAGGER_URL_REMOTE)

	// swagger route
	r.GET("/swagger/*any", ginSwagger.
		WrapHandler(
			swaggerFiles.Handler,
			conf,
		))

	return r
}
