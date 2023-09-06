package services

import (
	"calculations"
	"common"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"models"
	"net/http"
	"os"
	"reflect"
	"repositories"
	"sync"
	"time"
)

type StatsService struct {
	statsRepository *repositories.StatsRepository
}

func NewStatsService(statsRepository *repositories.StatsRepository) *StatsService {
	return &StatsService{
		statsRepository: statsRepository,
	}
}

func EmptyFieldErr(msg string) *models.ResponseError {
	return &models.ResponseError{
		Message: msg,
		Status:  http.StatusInternalServerError,
	}
}

func ValidateStatsInput(stats *[]models.Stats) (err *models.ResponseError) {
	var msg string
	for i, s := range *stats {
		statsVal := reflect.ValueOf(s)
		statsType := reflect.TypeOf(s)
		fieldsCount := statsVal.NumField()

		for j := 0; j < fieldsCount; j += 1 {
			fieldName := statsType.Field(j).Name
			val := statsVal.Field(j)
			isSet := val.IsValid() && !val.IsZero()

			if !isSet {
				msg = fmt.Sprintf("%v\nField \"%v\" for item with idx \"%v\" is not set", msg, fieldName, i)
			}
		}
	}
	if msg != "" {
		err = EmptyFieldErr(msg)
	}

	return err
}

func validateTimeFields(stats *models.Stats) *models.ResponseError {
	if stats.TimeFrame == "" {
		return &models.ResponseError{
			Message: "Invalid timeframe",
			Status:  http.StatusInternalServerError,
		}
	}
	if stats.Timestamp == 0 {
		return &models.ResponseError{
			Message: "Invalid stats timestamp",
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func processData(points []common.Point) *[]models.Stats {
	fStart := time.Now()
	// content, err := os.Open("../Archive/data.json")
	// common.ErrCheck(err)
	// defer content.Close()

	// var ds common.PointsSet
	// err2 := json.NewDecoder(content).Decode(&ds)
	// common.ErrCheck(err2)

	// binStorage, err3 := os.Create("../Archive/data.bin")
	// common.ErrCheck(err3)
	// defer binStorage.Close()

	// for _, p := range ds {
	// 	b := make([]byte, 16)
	// 	binary.LittleEndian.PutUint64(b[:8], p.Timestamp)
	// 	binary.LittleEndian.PutUint64(b[8:16], math.Float64bits(p.Rate))

	// 	_, wbErr := binStorage.Write(b)
	// 	common.ErrCheck(wbErr)
	// }

	binContent, err := os.Open("../Archive/data.bin")
	common.ErrCheck(err)
	defer binContent.Close()

	bytes, readErr := io.ReadAll(binContent)
	common.ErrCheck(readErr)

	dsFromBin := make(common.PointsSet, common.Hrs24pointsCount)
	gorutinesCount := 10
	dataChunkSize := len(dsFromBin) / gorutinesCount
	var wg sync.WaitGroup

	for i := dataChunkSize; i <= len(dsFromBin); i += dataChunkSize {
		innerLoopStart := i - dataChunkSize
		innerLoopEnd := i

		wg.Add(1)
		func(start, end int) {
			defer wg.Done()

			for j := start; j < end; j += 1 {
				offset := j * common.ByteChunkSize
				b := bytes[offset : offset+common.ByteChunkSize]

				ts := binary.LittleEndian.Uint64(b[:8])
				rateUint64 := binary.LittleEndian.Uint64(b[8:16])
				rate := math.Float64frombits(rateUint64)

				point := common.Point{
					Timestamp: ts,
					Rate:      rate,
				}
				dsFromBin[j] = point
			}
		}(innerLoopStart, innerLoopEnd)
	}
	wg.Wait()

	log.Printf("Reading + decoding took: %v", time.Since(fStart))
	fmt.Printf("ds size: %v\n\n", len(dsFromBin))

	startS := time.Now()
	newPoints := points
	remainedDs := dsFromBin[:len(dsFromBin)-len(points)]
	newPoints = append(newPoints, remainedDs...)
	log.Printf("Data set updating took: %v\n\n", time.Since(startS))

	startP := time.Now()
	parallelRes := calculations.Parallel(newPoints)
	log.Printf("Parallel calculations took: %v\n\n", time.Since(startP))
	var calculatedStats []models.Stats

	for _, mins5StatsItem := range parallelRes.Mins5 {
		mins5StatsItem.TimeFrame = "mins5"
		calculatedStats = append(calculatedStats, mins5StatsItem)
	}
	for _, mins30StatsItem := range parallelRes.Mins30 {
		mins30StatsItem.TimeFrame = "mins30"
		calculatedStats = append(calculatedStats, mins30StatsItem)
	}
	for _, hrs4StatsItem := range parallelRes.Hrs4 {
		hrs4StatsItem.TimeFrame = "hrs4"
		calculatedStats = append(calculatedStats, hrs4StatsItem)
	}
	for _, hrs24StatsItem := range parallelRes.Hrs24 {
		hrs24StatsItem.TimeFrame = "hrs24"
		calculatedStats = append(calculatedStats, hrs24StatsItem)
	}

	go func() {
		parallelOut, err3 := json.Marshal(parallelRes)
		common.ErrCheck(err3)
		os.WriteFile("./output/parallelOut-fromBin.json", parallelOut, 0644)
	}()

	return &calculatedStats
}

func (ss StatsService) SaveStats(points []common.Point) (*[]models.StoredStatsDB, *models.ResponseError) {
	stats := processData(points)

	validationErr := ValidateStatsInput(stats)
	if validationErr != nil {
		return nil, validationErr
	}

	return ss.statsRepository.SaveStats(stats)
}

func (ss StatsService) UpdateStatsRecord(dataToUpdate *models.Stats, recordId int) *models.ResponseError {
	validationErr := validateTimeFields(dataToUpdate)
	if validationErr != nil {
		return validationErr
	}
	t := time.Unix(int64(dataToUpdate.Timestamp/1000), 0).UTC()
	formattedTimestamp := t.Format(time.RFC3339)

	return ss.statsRepository.UpdateStatsRecord(&models.StoredStatsDB{
		Id:        recordId,
		Timestamp: formattedTimestamp,
		TimeFrame: dataToUpdate.TimeFrame,
	})
}

func (ss StatsService) GetStatsRecord(statsId int) (*models.StoredStatsDB, *models.ResponseError) {
	statsItem, err := ss.statsRepository.GetStatsOne(statsId)
	if err != nil {
		return nil, err
	}
	return statsItem, nil
}

func (ss StatsService) GetStatsByCreatedAt(creationTimestamp string) (*[]models.StoredStatsDB, *models.ResponseError) {
	statsItems, err := ss.statsRepository.GetStatsByCreatedAt(creationTimestamp)
	if err != nil {
		return nil, err
	}
	return statsItems, nil
}
