package parallel

import (
	"fmt"
	"sync"

	"common"
)

func stastLoopCalc(step, increment int, prevLvlStats []common.Stats, currentLvlStats *[]common.Stats) {
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

func Parallel(ds common.PointsSet) common.ResultsByTime {
	const mins5inMins30hrs4inHrs24 = 6
	const mins30inHrs4 = 8
	const mins5pointsCount = 5 * 60 * 100
	const mins30pointsCount = mins5pointsCount * mins5inMins30hrs4inHrs24
	const hrs4pointsCount = mins30pointsCount * mins30inHrs4
	const hrs24pointsCount = hrs4pointsCount * mins5inMins30hrs4inHrs24

	mins5rangeSize, dsLoopIncrement := common.CalcSize(len(ds), mins5pointsCount, mins5pointsCount)
	mins30RangeSize, mins5loopIncrement := common.CalcSize(len(ds), mins30pointsCount, mins5inMins30hrs4inHrs24)
	hrs4rangeSize, mins30loopIncrement := common.CalcSize(len(ds), hrs4pointsCount, mins30inHrs4)
	hrs24rangeSize, hrs4loopIncrement := common.CalcSize(len(ds), hrs24pointsCount, mins5inMins30hrs4inHrs24)

	var mins5statsRange = make([]common.Stats, mins5rangeSize)
	var min30statsRange = make([]common.Stats, mins30RangeSize)
	var hrs4statsRange = make([]common.Stats, hrs4rangeSize)
	var hrs24statsRange = make([]common.Stats, hrs24rangeSize)

	count5minsStats := 0

	var wg5mins sync.WaitGroup

	for i := mins5pointsCount; i <= len(ds)+dsLoopIncrement; i += mins5pointsCount {
		wg5mins.Add(1)

		go func(pointsIdx, stats5minsIdx int, currentLvlStats *[]common.Stats) {
			defer wg5mins.Done()

			end := pointsIdx
			if pointsIdx > len(ds) {
				end = len(ds)
			}
			points := ds[pointsIdx-mins5pointsCount : end]
			current5minsStats := points.CalcPoints()
			(*currentLvlStats)[stats5minsIdx] = current5minsStats

		}(i, count5minsStats, &mins5statsRange)

		count5minsStats += 1
	}
	wg5mins.Wait()

	stastLoopCalc(mins5inMins30hrs4inHrs24, mins5loopIncrement, mins5statsRange, &min30statsRange)
	stastLoopCalc(mins30inHrs4, mins30loopIncrement, min30statsRange, &hrs4statsRange)
	stastLoopCalc(mins5inMins30hrs4inHrs24, hrs4loopIncrement, hrs4statsRange, &hrs24statsRange)

	for _, v := range hrs24statsRange[:] {
		fmt.Printf("Avg: %v; High: %v; Low: %v; Open: %v; Close: %v\n", v.Average, v.High, v.Low, v.Open, v.Close)
	}

	return common.ResultsByTime{
		Mins5:  mins5statsRange,
		Mins30: min30statsRange,
		Hrs4:   hrs4statsRange,
		Hrs24:  hrs24statsRange,
	}
}
