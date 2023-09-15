package common

import (
	"models"
	"time"
)

type Point struct {
	Rate      float64
	Timestamp uint64
}

type PointsSet []Point

type StatsSet []models.Stats

func ErrCheck(e error) {
	if e != nil {
		panic(e)
	}
}

const ROLE_ADMIN = "admin"
const ROLE_USER = "user"

const ByteChunkSize = 16
const Mins5inMins30hrs4inHrs24 = 6
const Mins30inHrs4 = 8
const Mins5pointsCount = 5 * 60 * 100
const Mins30pointsCount = Mins5pointsCount * Mins5inMins30hrs4inHrs24
const Hrs4pointsCount = Mins30pointsCount * Mins30inHrs4
const Hrs24pointsCount = Hrs4pointsCount * Mins5inMins30hrs4inHrs24

func FormatTimestamp(numericTst uint64) string {
	t := time.Unix(int64(numericTst/1000), 0).UTC()
	formattedTimestamp := t.Format(time.RFC3339)

	return formattedTimestamp
}

func (ds PointsSet) CalcPoints() (stats models.Stats) {
	var total float64

	stats = models.Stats{
		Timestamp: ds[len(ds)-1].Timestamp,
		Average:   0,
		Open:      ds[len(ds)-1].Rate,
		Close:     ds[0].Rate,
		High:      ds[0].Rate,
		Low:       ds[0].Rate,
	}
	for _, v := range ds {
		total += v.Rate
		if stats.High < v.Rate {
			stats.High = v.Rate
		}
		if stats.Low > v.Rate {
			stats.Low = v.Rate
		}
	}
	// slices.SortFunc[Set, Point](ds, func(a, b Point) int {
	// 	return int(a.Rate - b.Rate)
	// })
	// stats.High = ds[len(ds) - 1].Rate
	// stats.Low = ds[0].Rate

	stats.Average = total / float64(len(ds))

	return
}

func (ds StatsSet) CalcStats() (stats models.Stats) {
	var total float64

	stats = models.Stats{
		Timestamp: ds[len(ds)-1].Timestamp,
		Average:   0,
		High:      ds[0].High,
		Low:       ds[0].Low,
		Open:      ds[len(ds)-1].Open,
		Close:     ds[0].Close,
	}
	for _, v := range ds {
		total += v.Average
		if stats.High < v.High {
			stats.High = v.High
		}
		if stats.Low > v.Low {
			stats.Low = v.Low
		}
		// stats.High = math.Max(stats.High, v.High)
		// stats.Low = math.Min(stats.Low, v.Low)
	}
	stats.Average = total / float64(len(ds))

	return
}

func CalcSize(dataSize, partSize, loopStep int) (int, int) {
	setsRangeSize := dataSize / partSize

	if dataSize%partSize == 0 {
		return setsRangeSize, 0
	}
	return setsRangeSize + 1, loopStep
}
