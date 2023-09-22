package services

import (
	"models"
	"net/http"
	"repositories"
	"strconv"
	"utils"
)

type NewUserParams struct {
	Name        string `json:"name" required:"true"`
	Password    string `json:"password" required:"true"`
	Role        string `json:"role" required:"true"`
	AccessToken string `json:"-"`
}

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(usersRepository *repositories.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

// func generateAccessToken(name string, id int) (string, *models.ResponseError) {
// 	baseStr := fmt.Sprintf("%v-%d", name, id)
// 	hash, err := bcrypt.GenerateFromPassword([]byte(baseStr), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", &models.ResponseError{
// 			Message: err.Error(),
// 			Status:  http.StatusInternalServerError,
// 		}
// 	}

// 	return base64.StdEncoding.EncodeToString(hash), nil
// }

func (us UsersService) AddUser(params NewUserParams) *models.ResponseError {
	name := params.Name
	password := params.Password
	role := params.Role
	accessToken := params.AccessToken

	if name == "" || password == "" || role == "" {
		return &models.ResponseError{
			Message: "Name, password and role are required",
			Status:  http.StatusBadRequest,
		}
	}
	if accessToken == "" && role == "admin" {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusUnauthorized,
		}
	}
	if accessToken == "" && role != "admin" {
		err := us.usersRepository.AddUser(name, password, role)
		if err != nil {
			return err
		}
		return nil
	}

	_, isAuthorized, _, err := us.AuthorizeUser(accessToken, []string{role})
	if err != nil {
		return err
	}

	if !isAuthorized {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}

func (us UsersService) Login(name, password string) (string, *models.ResponseError) {
	if name == "" || password == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}
	id, role, err := us.usersRepository.LoginUser(name, password)
	if err != nil {
		return "", err
	}
	if id == 0 {
		return "", &models.ResponseError{
			Message: "Login failed",
			Status:  http.StatusUnauthorized,
		}
	}
	accessToken, err := utils.GenerateJWT(name, role, id) // generateAccessToken(name, id)
	if err != nil {
		return "", err
	}

	us.usersRepository.SetAccessToken(accessToken, id)
	return accessToken, nil
}

func (us UsersService) Logout(userId int) *models.ResponseError {
	return us.usersRepository.RemoveAccessToken(userId)
}

func (us UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (int, bool, string, *models.ResponseError) {
	_, _, err := us.usersRepository.GetUser(accessToken)
	if err != nil {
		return 0, false, "", err
	}

	tokenClaims, err := utils.VerifyJWT(accessToken)
	if err != nil {
		return 0, false, "", err
	}
	userId, parsingErr := strconv.ParseInt(tokenClaims.StandardClaims.Subject, 0, 0)
	if parsingErr != nil {
		return 0, false, "", err
	}
	role := tokenClaims.Role

	if err != nil {
		return 0, false, "", err
	}
	if role == "" {
		return 0, false, "", &models.ResponseError{
			Message: "Failed to authorize user",
			Status:  http.StatusUnauthorized,
		}
	}

	for _, expected := range expectedRoles {
		if expected == role {
			return int(userId), true, role, nil
		}
	}

	return int(userId), false, role, nil
}
