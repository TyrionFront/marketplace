package controllers

import (
	"common"
	"encoding/json"
	"io"
	"log"
	"models"
	"net/http"
	"services"
	"strconv"

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

	userId, isAuth, _, authErr := sc.usersService.AuthorizeUser(accessToken, expectedRoles)
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
	if len(reqBody) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ResponseError{
			Message: "Request body can't be empty.",
			Status:  http.StatusBadRequest,
		})
		return
	}
	var points []common.Point
	err = json.Unmarshal(reqBody, &points)
	if err != nil {
		log.Println("Error while unmarshaling request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(points) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ResponseError{
			Message: "No new data for recalculation.",
			Status:  http.StatusBadRequest,
		})
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

func (sc StatsController) PrepareStatsByUser(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	expectedRoles := []string{common.ROLE_ADMIN, common.ROLE_USER}

	tokenUserId, isAuth, role, authErr := sc.usersService.AuthorizeUser(accessToken, expectedRoles)
	if authErr != nil {
		ctx.JSON(authErr.Status, authErr)
		return
	}
	if !isAuth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	userIdStr := ctx.Param("userId")
	parsedUserId, err := strconv.ParseInt(userIdStr, 0, 0)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	userId := int(parsedUserId)

	if role != common.ROLE_ADMIN && tokenUserId != userId {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ResponseError{
			Message: "Unauthorized request",
			Status:  http.StatusBadRequest,
		})
		return
	}

	res, resError := sc.statsService.GetStatsByUser(userId)
	if resError != nil {
		ctx.AbortWithStatusJSON(resError.Status, resError)
		return
	}

	if len(*res) == 0 {
		ctx.JSON(http.StatusOK, "No data has been found for the current user")
		return
	}
	ctx.JSON(http.StatusOK, models.StatsByUser{
		Size: len(*res),
		Data: res,
	})
}
