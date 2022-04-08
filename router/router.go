package router

import (
	"net/http"

	"github.com/asrofilfachrulr/get-nearby-places/controller"
	"github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// const SWAGGER_URL_REMOTE = "https://gist.githubusercontent.com/yohang88/2efb1f26f452d059643fb7ea00c15a10/raw/3b775398f0365d85ea6eb200f8f192091f1fdef1/jcc-openapi-spec.yaml"

const SWAGGER_LOCAL = "./swagger-config.yaml"

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
	conf := ginSwagger.URL(SWAGGER_LOCAL)

	// swagger route
	r.GET("/swagger/*any", ginSwagger.
		WrapHandler(
			swaggerFiles.Handler,
			conf,
		))

	return r
}
