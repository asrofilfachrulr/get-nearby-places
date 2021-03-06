package controller

import (
	"net/http"
	"strconv"

	"github.com/asrofilfachrulr/get-nearby-places/models"
	"github.com/gin-gonic/gin"
)

func GetNearby(ctx *gin.Context) {
	bplace := ctx.MustGet("places").(models.BatchPlace)

	// error ignored, assuming user input is always valid input
	lat, _ := strconv.ParseFloat(ctx.Query("latitude"), 64)
	lon, _ := strconv.ParseFloat(ctx.Query("longitude"), 64)
	cid, _ := strconv.ParseUint(ctx.Query("category_id"), 10, 8)

	q := models.WebQuery{
		Latitude:   lat,
		Longitude:  lon,
		CategoryId: uint8(cid),
	}

	result, _ := models.GetNearbyPlaces(q, bplace)

	ctx.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": len(result),
	})
}
