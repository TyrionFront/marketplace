package repositories

import (
	"database/sql"
	"models"
	"net/http"
	"time"
)

type UsersRepository struct {
	dbHandler *sql.DB
}

func NewUsersRepository(dbHandler *sql.DB) *UsersRepository {
	return &UsersRepository{
		dbHandler: dbHandler,
	}
}

func (ur UsersRepository) AddUser(name, password, role string) *models.ResponseError {
	query := `
		INSERT INTO users(username, user_password, user_role)
		VALUES($1, crypt($2, gen_salt('bf')), $3)
	`
	_, err := ur.dbHandler.Exec(query, name, password, role)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur UsersRepository) LoginUser(name, password string) (int, string, *models.ResponseError) {
	query := `
		SELECT id, user_role
		FROM users
		WHERE username = $1 AND user_password = crypt($2, user_password)
	`

	rows, err := ur.dbHandler.Query(query, name, password)
	if err != nil {
		return 0, "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var id int
	var userRole string
	for rows.Next() {
		err := rows.Scan(&id, &userRole)
		if err != nil {
			return 0, "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if err = rows.Err(); err != nil {
		return 0, "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return id, userRole, nil
}

func (ur UsersRepository) GetUser(accessToken string) (int, string, *models.ResponseError) {
	query := `
		SELECT id, user_role
		FROM users
		WHERE access_token = $1
	`

	rows, err := ur.dbHandler.Query(query, accessToken)
	if err != nil {
		return 0, "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var id int
	var role, expiresAt string
	for rows.Next() {
		err := rows.Scan(&id, &role, &expiresAt)
		if err != nil {
			return 0, "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if err = rows.Err(); err != nil {
		return 0, "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if id == 0 || expiresAt == "" {
		return 0, "", &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	parsedTimePassed, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return 0, "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	unixTimePassed := parsedTimePassed.Unix()

	if time.Now().Unix() > unixTimePassed {
		return 0, "", &models.ResponseError{
			Message: "Expired token. Please log in again",
			Status:  http.StatusUnauthorized,
		}
	}
	ur.SetAccessToken(accessToken, id)

	return id, role, nil
}

func (ur UsersRepository) SetAccessToken(accessToken string, id int) *models.ResponseError {
	query := `
		UPDATE users
		SET access_token = $1
		WHERE id = $2
	`

	_, err := ur.dbHandler.Exec(query, accessToken, id)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur UsersRepository) RemoveAccessToken(accessToken string) *models.ResponseError {
	query := `
		UPDATE users
		SET access_token = ''
		WHERE id = $1
	`
	_, err := ur.dbHandler.Exec(query, accessToken)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
