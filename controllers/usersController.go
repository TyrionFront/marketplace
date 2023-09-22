package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"services"
	"utils"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	usersService *services.UsersService
}

func NewUsersController(usersService *services.UsersService) *UsersController {
	return &UsersController{
		usersService: usersService,
	}
}

func (uc UsersController) AddUser(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var newUserParams services.NewUserParams
	err = json.Unmarshal(reqBody, &newUserParams)
	if err != nil {
		log.Println("Error while unmarshaling request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var accessToken string
	if newUserParams.Role == "admin" {
		token, extrErr := utils.ExtractToken(ctx)
		if extrErr != nil {
			ctx.JSON(extrErr.Status, extrErr)
			return
		}
		accessToken = token
	}
	newUserParams.AccessToken = accessToken

	addingErr := uc.usersService.AddUser(newUserParams)
	if addingErr != nil {
		ctx.AbortWithStatusJSON(addingErr.Status, addingErr)
		return
	}

	successMsg := "User has been successfully added.\nPlease log in using your username and password"
	ctx.JSON(http.StatusCreated, successMsg)
}

func (uc UsersController) Login(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		log.Println("Error while reading credentials")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accessToken, err := uc.usersService.Login(username, password)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Status, err)
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

func (uc UsersController) Logout(ctx *gin.Context) {
	accessToken, extrErr := utils.ExtractToken(ctx)
	if extrErr != nil {
		ctx.JSON(extrErr.Status, extrErr)
		return
	}
	err := uc.usersService.Logout(accessToken)

	if err != nil {
		ctx.AbortWithStatusJSON(err.Status, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
