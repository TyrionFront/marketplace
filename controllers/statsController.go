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
	statsService *services.StatsService
	usersService *services.UsersService
}

func NewStatsController(statsService *services.StatsService, usersService *services.UsersService) *StatsController {
	return &StatsController{
		statsService: statsService,
		usersService: usersService,
	}
}

func (sc StatsController) SaveStats(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	expectedRoles := []string{common.ROLE_ADMIN, common.ROLE_USER}

	userId, isAuth, authErr := sc.usersService.AuthorizeUser(accessToken, expectedRoles)
	if authErr != nil {
		ctx.JSON(authErr.Status, authErr)
		return
	}
	if !isAuth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	if userId == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
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

	res, resError := sc.statsService.SaveStats(points, userId)
	if resError != nil {
		ctx.AbortWithStatusJSON(resError.Status, resError)
		return
	}
	if len(*res) == 0 {
		ctx.JSON(http.StatusOK, "No new data has been saved")
		return
	}
	ctx.JSON(http.StatusCreated, res)
}
