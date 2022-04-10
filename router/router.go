package router

import (
	"net/http"

	"github.com/asrofilfachrulr/get-nearby-places/controller"
	"github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const SWAGGER_URL_REMOTE = "https://gist.githubusercontent.com/asrofilfachrulr/1c8d2a3bf151678c8e01071facb553ed/raw/227740ee0b51eb3d3e95b0512840d616aa65bcf5/get-nearby-places.yaml"

// const SWAGGER_LOCAL = "./swagger-config.yaml"

func SetupRouter(data models.BatchPlace) *gin.Engine {
	r := gin.Default()

	// attach BatchPlace to gin Context
	r.Use(func(ctx *gin.Context) {
		ctx.Set("places", data)
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "mantap")
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
