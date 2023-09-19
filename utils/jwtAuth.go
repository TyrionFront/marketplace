package utils

import (
	"fmt"
	"models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	UserName string `json:"username" required:"true"`
	Role     string `json:"role" required:"true"`
	jwt.StandardClaims
}

func ExtractToken(ctx *gin.Context) (string, *models.ResponseError) {
	bearerToken := ctx.Request.Header.Get("Authorization")
	if bearerToken == "" {
		return "", &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	tokenParts := strings.Split(bearerToken, " ")
	if len(tokenParts) != 2 {
		return "", &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	return tokenParts[1], nil
}

func GenerateJWT(userName, role string, userId int) (string, *models.ResponseError) {
	_, checkErr := os.Stat("./keys/private.pem")
	if os.IsNotExist(checkErr) {
		fmt.Println(checkErr)
		GenerateKeyPair()
	}

	privateKeyBytes, err := os.ReadFile("./keys/private.pem")
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	jwtTTLStr, envErr := GetEnvValue("JWT_TTL_SEC")
	if envErr != nil {
		return "", envErr
	}
	jwtTTL, err := strconv.ParseInt(jwtTTLStr, 0, 0)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	claims := JWTClaims{
		UserName: userName,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			Subject:   fmt.Sprint(userId),
			ExpiresAt: time.Now().Add(time.Duration(jwtTTL) * time.Second).Unix(),
		},
	}

	initialToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, signingErr := initialToken.SignedString(privateKeyBytes)
	if signingErr != nil {
		return "", &models.ResponseError{
			Message: signingErr.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return signedToken, nil
}

func VerifyJWT(token string) (JWTClaims, *models.ResponseError) {
	privateKeyBytes, err := os.ReadFile("./keys/private.pem")
	if err != nil {
		return JWTClaims{}, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(tk *jwt.Token) (interface{}, error) {
		return privateKeyBytes, nil
	})
	if err != nil {
		return JWTClaims{}, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if !parsedToken.Valid {
		return JWTClaims{}, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		return JWTClaims{}, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	return *claims, nil
}
