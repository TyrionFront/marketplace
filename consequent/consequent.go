package consequent

import (
	"fmt"

	"common"
)

func Consequent(ds common.PointsSet) common.ResultsByTime {
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
	count30minsStats := 0
	count4hrsStats := 0
	count24hrsStats := 0

	for i := mins5pointsCount; i <= len(ds); i += mins5pointsCount {
		if i%mins5pointsCount == 0 {
			pointsToCalc := ds[i-mins5pointsCount : i]
			current5minsStats := pointsToCalc.CalcPoints()

			mins5statsRange[count5minsStats] = current5minsStats
			count5minsStats += 1
		}
		if i%mins30pointsCount == 0 {
			start := count5minsStats - mins5inMins30hrs4inHrs24
			var statsToCalc common.StatsSet = mins5statsRange[start:count5minsStats]

			current30minsStats := statsToCalc.CalcStats()
			min30statsRange[count30minsStats] = current30minsStats
			count30minsStats += 1
		}
		if i%hrs4pointsCount == 0 {
			start := count30minsStats - mins30inHrs4
			var statsToCalc common.StatsSet = min30statsRange[start:count30minsStats]

			current4hrsStats := statsToCalc.CalcStats()
			hrs4statsRange[count4hrsStats] = current4hrsStats
			count4hrsStats += 1
		}
		if i%hrs24pointsCount == 0 {
			start := count4hrsStats - mins5inMins30hrs4inHrs24
			var statsToCalc common.StatsSet = hrs4statsRange[start:count4hrsStats]

			current4hrsStats := statsToCalc.CalcStats()
			hrs24statsRange[count24hrsStats] = current4hrsStats
			count24hrsStats += 1
		}
	}

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
