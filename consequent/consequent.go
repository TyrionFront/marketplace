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

	mins5rangeSize, dsLoopIncrement := common.CalcSize(len(ds), mins5pointsCount, mins5pointsCount)
	mins30RangeSize, mins5loopIncrement := common.CalcSize(len(ds), mins30pointsCount, mins5inMins30hrs4inHrs24)
	hrs4rangeSize, mins30loopIncrement := common.CalcSize(len(ds), hrs4pointsCount, mins30inHrs4)
	hrs24rangeSize, hrs4loopIncrement := common.CalcSize(len(ds), hrs24pointsCount, mins5inMins30hrs4inHrs24)

	var mins5statsRange = make([]common.Stats, mins5rangeSize)
	var min30statsRange = make([]common.Stats, mins30RangeSize)
	var hrs4statsRange = make([]common.Stats, hrs4rangeSize)
	var hrs24statsRange = make([]common.Stats, hrs24rangeSize)

	count5minsStats := 0
	count30minsStats := 0
	count4hrsStats := 0
	count24hrsStats := 0

	for i := mins5pointsCount; i <= len(ds)+dsLoopIncrement; i += mins5pointsCount {
		end := i
		if i > len(ds) {
			end = len(ds)
		}
		pointsToCalc := ds[i-mins5pointsCount : end]
		current5minsStats := pointsToCalc.CalcPoints()

		mins5statsRange[count5minsStats] = current5minsStats
		count5minsStats += 1
	}

	for i := mins5inMins30hrs4inHrs24; i <= len(mins5statsRange)+mins5loopIncrement; i += mins5inMins30hrs4inHrs24 {
		end := i
		if i > len(mins5statsRange) {
			end = len(mins5statsRange)
		}
		var statsToCalc common.StatsSet = mins5statsRange[i-mins5inMins30hrs4inHrs24 : end]

		current30minsStats := statsToCalc.CalcStats()
		min30statsRange[count30minsStats] = current30minsStats
		count30minsStats += 1
	}

	for i := mins30inHrs4; i <= len(min30statsRange)+mins30loopIncrement; i += mins30inHrs4 {
		end := i
		if i > len(min30statsRange) {
			end = len(min30statsRange)
		}
		var statsToCalc common.StatsSet = min30statsRange[i-mins30inHrs4 : end]

		current4hrsStats := statsToCalc.CalcStats()
		hrs4statsRange[count4hrsStats] = current4hrsStats
		count4hrsStats += 1
	}

	for i := mins5inMins30hrs4inHrs24; i <= len(hrs4statsRange)+hrs4loopIncrement; i += mins5inMins30hrs4inHrs24 {
		end := i
		if i > len(hrs4statsRange) {
			end = len(hrs4statsRange)
		}
		var statsToCalc common.StatsSet = hrs4statsRange[i-mins5inMins30hrs4inHrs24 : end]

		current4hrsStats := statsToCalc.CalcStats()
		hrs24statsRange[count24hrsStats] = current4hrsStats
		count24hrsStats += 1
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
