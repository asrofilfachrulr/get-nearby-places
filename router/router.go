package router

import (
	"github.com/asrofilfachrulr/get-nearby-places/controller"
	"github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(data []models.Place) *gin.Engine {
	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		ctx.Set("places", data)
	})

	r.GET("/search", controller.GetNearby)

	conf := ginSwagger.URL("https://gist.githubusercontent.com/yohang88/2efb1f26f452d059643fb7ea00c15a10/raw/3b775398f0365d85ea6eb200f8f192091f1fdef1/jcc-openapi-spec.yaml")

	// swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, conf))

	return r
}