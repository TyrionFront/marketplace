package parallel

import (
	"fmt"
	"sync"

	"common"
)

func stastLoopCalc(step int, prevLvlStats, currentLvlStats []common.Stats) {
	var wg sync.WaitGroup
	curentStatsCount := 0

	for i := step; i <= len(prevLvlStats); i += step {
		wg.Add(1)

		go func(lowerStatsIdx, higherStatsIdx int) {
			defer wg.Done()

			start := lowerStatsIdx - step
			var statsToCalc common.StatsSet = prevLvlStats[start:lowerStatsIdx]

			current30minsStats := statsToCalc.CalcStats()
			currentLvlStats[higherStatsIdx] = current30minsStats
		}(i, curentStatsCount)

		curentStatsCount += 1
	}
	wg.Wait()
}

func Parallel(ds common.PointsSet) {
	const mins5inMins30hrs4inHrs24 = 6
	const mins30inHrs4 = 8

	const mins5pointsCount = 5 * 60 * 100
	const mins30pointsCount = mins5pointsCount * mins5inMins30hrs4inHrs24
	const hrs4pointsCount = mins30pointsCount * mins30inHrs4
	const hrs24pointsCount = hrs4pointsCount * mins5inMins30hrs4inHrs24

	var mins5statsRange = make([]common.Stats, len(ds)/mins5pointsCount)
	var min30statsRange = make([]common.Stats, len(ds)/mins30pointsCount)
	var hrs4statsRange = make([]common.Stats, len(ds)/hrs4pointsCount)
	var hrs24statsRange = make([]common.Stats, len(ds)/hrs24pointsCount)

	count5minsStats := 0

	var wg5mins sync.WaitGroup

	for i := mins5pointsCount; i <= len(ds); i += mins5pointsCount {
		wg5mins.Add(1)

		go func(pointsIdx, stats5minsIdx int) {
			defer wg5mins.Done()

			points := ds[pointsIdx-mins5pointsCount : pointsIdx]
			current5minsStats := points.CalcPoints()
			mins5statsRange[stats5minsIdx] = current5minsStats
		}(i, count5minsStats)

		count5minsStats += 1
	}
	wg5mins.Wait()

	stastLoopCalc(mins5inMins30hrs4inHrs24, mins5statsRange, min30statsRange)
	stastLoopCalc(mins30inHrs4, min30statsRange, hrs4statsRange)
	stastLoopCalc(mins5inMins30hrs4inHrs24, hrs4statsRange, hrs24statsRange)

	for _, v := range hrs24statsRange[:] {
		fmt.Printf("Avg: %v; High: %v; Low: %v; Open: %v; Close: %v\n", v.Average, v.High, v.Low, v.Open, v.Close)
	}
}
