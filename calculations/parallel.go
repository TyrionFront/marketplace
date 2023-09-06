package calculations

import (
	"sync"

	"common"
	"models"
)

func stastLoopCalc(step, increment int, prevLvlStats []models.Stats, currentLvlStats *[]models.Stats) {
	var wg sync.WaitGroup
	curentStatsCount := 0

	for i := step; i <= len(prevLvlStats)+increment; i += step {
		wg.Add(1)

		end := i
		if end > len(prevLvlStats) {
			end = len(prevLvlStats)
		}
		go func(startIdx, endIdx, higherLvlStatsIdx int) {
			defer wg.Done()

			var statsToCalc common.StatsSet = prevLvlStats[startIdx:endIdx]

			currentStats := statsToCalc.CalcStats()
			(*currentLvlStats)[higherLvlStatsIdx] = currentStats

		}(i-step, end, curentStatsCount)

		curentStatsCount += 1
	}
	wg.Wait()
}

func Parallel(ds common.PointsSet) models.ResultsByTime {
	mins5rangeSize, dsLoopIncrement := common.CalcSize(len(ds), common.Mins5pointsCount, common.Mins5pointsCount)
	mins30RangeSize, mins5loopIncrement := common.CalcSize(len(ds), common.Mins30pointsCount, common.Mins5inMins30hrs4inHrs24)
	hrs4rangeSize, mins30loopIncrement := common.CalcSize(len(ds), common.Hrs4pointsCount, common.Mins30inHrs4)
	hrs24rangeSize, hrs4loopIncrement := common.CalcSize(len(ds), common.Hrs24pointsCount, common.Mins5inMins30hrs4inHrs24)

	var mins5statsRange = make([]models.Stats, mins5rangeSize)
	var min30statsRange = make([]models.Stats, mins30RangeSize)
	var hrs4statsRange = make([]models.Stats, hrs4rangeSize)
	var hrs24statsRange = make([]models.Stats, hrs24rangeSize)

	count5minsStats := 0

	var wg5mins sync.WaitGroup

	for i := common.Mins5pointsCount; i <= len(ds)+dsLoopIncrement; i += common.Mins5pointsCount {
		wg5mins.Add(1)

		go func(pointsIdx, stats5minsIdx int, currentLvlStats *[]models.Stats) {
			defer wg5mins.Done()

			end := pointsIdx
			if pointsIdx > len(ds) {
				end = len(ds)
			}
			points := ds[pointsIdx-common.Mins5pointsCount : end]
			current5minsStats := points.CalcPoints()
			(*currentLvlStats)[stats5minsIdx] = current5minsStats

		}(i, count5minsStats, &mins5statsRange)

		count5minsStats += 1
	}
	wg5mins.Wait()

	stastLoopCalc(common.Mins5inMins30hrs4inHrs24, mins5loopIncrement, mins5statsRange, &min30statsRange)
	stastLoopCalc(common.Mins30inHrs4, mins30loopIncrement, min30statsRange, &hrs4statsRange)
	stastLoopCalc(common.Mins5inMins30hrs4inHrs24, hrs4loopIncrement, hrs4statsRange, &hrs24statsRange)

	return models.ResultsByTime{
		Mins5:  mins5statsRange,
		Mins30: min30statsRange,
		Hrs4:   hrs4statsRange,
		Hrs24:  hrs24statsRange,
	}
}
