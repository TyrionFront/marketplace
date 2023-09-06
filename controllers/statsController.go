package controllers

import (
	"common"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"services"

	"github.com/gin-gonic/gin"
)

type StatsController struct {
	statsService services.StatsService
}

func NewStatsController(statsService *services.StatsService) *StatsController {
	return &StatsController{
		statsService: *statsService,
	}
}

func (sc StatsController) SaveStats(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		log.Println("Error while reading request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var points []common.Point
	err = json.Unmarshal(reqBody, &points)
	if err != nil {
		log.Println("Error while unmarshaling request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res, resError := sc.statsService.SaveStats(points)
	if resError != nil {
		ctx.AbortWithStatusJSON(resError.Status, resError)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}
