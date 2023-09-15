package repositories

import (
	"common"
	"database/sql"
	"fmt"
	"models"
	"net/http"
	"reflect"
	"time"
)

type StatsRepository struct {
	dbHandler *sql.DB
	// transaction *sql.Tx
}

func NewStatsRepository(dbHandler *sql.DB) *StatsRepository {
	return &StatsRepository{
		dbHandler: dbHandler,
	}
}

func mapResults(rows *sql.Rows) (*[]models.StoredStatsDB, *models.ResponseError) {
	var storedStatsDB []models.StoredStatsDB
	var id, user int
	var createdAt, updatedAt, timeFrame, timestampDB string
	var average, high, low, open, close float64

	var err error
	for rows.Next() {
		err = rows.Scan(
			&id, &createdAt, &updatedAt, &timestampDB,
			&average, &high, &low, &open, &close, &timeFrame, &user,
		)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		statsRecord := models.StoredStatsDB{
			Id:          id,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Timestamp:   timestampDB,
			TimeFrame:   timeFrame,
			Average:     average,
			High:        high,
			Low:         low,
			Open:        open,
			Close:       close,
			RelatedUser: user,
		}
		storedStatsDB = append(storedStatsDB, statsRecord)
	}
	if err = rows.Err(); err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &storedStatsDB, nil
}

func (sr StatsRepository) SaveStats(stats *[]models.Stats, user int) (*[]models.StoredStatsDB, *models.ResponseError) {
	query := `
		INSERT INTO calculated_stats(timestamp, average, high, low, open, close, time_frame, created_at, updated_at, related_user)
		VALUES 
	`
	queryParams := []interface{}{}
	totalCount := 0
	for i, s := range *stats {
		fieldsCount := reflect.ValueOf(s).NumField()

		query += `(`
		for j := 0; j < fieldsCount+3; j += 1 {
			totalCount += 1
			var paramsSetEnd string
			if j < fieldsCount+2 {
				paramsSetEnd = `, `
			}
			query += `$` + fmt.Sprint(totalCount, paramsSetEnd)
		}
		paramsPartEnd := `), `
		if i == len(*stats)-1 {
			paramsPartEnd = `)`
		}
		query += paramsPartEnd

		formattedTimestamp := common.FormatTimestamp(s.Timestamp)
		currentTimeStamp := time.Now().Format(time.RFC3339)
		queryParams = append(
			queryParams, formattedTimestamp, s.Average, s.High,
			s.Low, s.Open, s.Close, s.TimeFrame, currentTimeStamp, currentTimeStamp, user,
		)
	}
	query += `
		RETURNING *
	`

	rows, err := sr.dbHandler.Query(
		query, queryParams...,
	)
	if err != nil {
		fmt.Println(query)
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	return mapResults(rows)
}

func (sr StatsRepository) UpdateStatsRecord(stats *models.StoredStatsDB) *models.ResponseError {
	query := `
		UPDATE calculated_stats
		SET
			timestamp = $2,
			time_frame = $3
		WHERE id = $1
	`
	res, err := sr.dbHandler.Exec(query, stats.Id, stats.Timestamp, stats.TimeFrame)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Stats record not found",
			Status:  http.StatusNotFound,
		}
	}
	return nil
}

func (sr StatsRepository) GetStatsOne(statsId int) (*models.StoredStatsDB, *models.ResponseError) {
	query := `
		SELECT *
		FROM calculated_stats
		WHERE id = $1
	`
	rows, err := sr.dbHandler.Query(query, statsId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var id int
	var createdAt, updatedAt, timeFrame, timestampDB string
	var average, high, low, open, close float64

	for rows.Next() {
		err := rows.Scan(
			&id, &createdAt, &updatedAt, &timestampDB,
			&average, &high, &low, &open, &close, &timeFrame,
		)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.StoredStatsDB{
		Id:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Timestamp: timestampDB,
		TimeFrame: timeFrame,
		Average:   average,
		High:      high,
		Low:       low,
		Open:      open,
		Close:     close,
	}, nil
}

func (sr StatsRepository) GetStatsByCreatedAt(creationTimestamp string) (*[]models.StoredStatsDB, *models.ResponseError) {
	query := `
		SELECT *
		FROM calculated_stats
		WHERE created_at = $1
	`

	rows, err := sr.dbHandler.Query(query, creationTimestamp)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	return mapResults(rows)
}

func (sr StatsRepository) GetAllStatsByUser(userId int) (*[]models.StoredStatsDB, *models.ResponseError) {
	query := `
		SELECT *
		FROM calculated_stats
		WHERE related_user = $1
	`

	rows, err := sr.dbHandler.Query(query, userId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	return mapResults(rows)
}
