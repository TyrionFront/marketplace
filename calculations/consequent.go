package calculations

import (
	"models"

	"common"
)

func Consequent(ds common.PointsSet) models.ResultsByTime {
	mins5rangeSize, dsLoopIncrement := common.CalcSize(len(ds), common.Mins5pointsCount, common.Mins5pointsCount)
	mins30RangeSize, mins5loopIncrement := common.CalcSize(len(ds), common.Mins30pointsCount, common.Mins5inMins30hrs4inHrs24)
	hrs4rangeSize, mins30loopIncrement := common.CalcSize(len(ds), common.Hrs4pointsCount, common.Mins30inHrs4)
	hrs24rangeSize, hrs4loopIncrement := common.CalcSize(len(ds), common.Hrs24pointsCount, common.Mins5inMins30hrs4inHrs24)

	var mins5statsRange = make([]models.Stats, mins5rangeSize)
	var min30statsRange = make([]models.Stats, mins30RangeSize)
	var hrs4statsRange = make([]models.Stats, hrs4rangeSize)
	var hrs24statsRange = make([]models.Stats, hrs24rangeSize)

	count5minsStats := 0
	count30minsStats := 0
	count4hrsStats := 0
	count24hrsStats := 0

	for i := common.Mins5pointsCount; i <= len(ds)+dsLoopIncrement; i += common.Mins5pointsCount {
		end := i
		if i > len(ds) {
			end = len(ds)
		}
		pointsToCalc := ds[i-common.Mins5pointsCount : end]
		current5minsStats := pointsToCalc.CalcPoints()

		mins5statsRange[count5minsStats] = current5minsStats
		count5minsStats += 1
	}

	for i := common.Mins5inMins30hrs4inHrs24; i <= len(mins5statsRange)+mins5loopIncrement; i += common.Mins5inMins30hrs4inHrs24 {
		end := i
		if i > len(mins5statsRange) {
			end = len(mins5statsRange)
		}
		var statsToCalc common.StatsSet = mins5statsRange[i-common.Mins5inMins30hrs4inHrs24 : end]

		current30minsStats := statsToCalc.CalcStats()
		min30statsRange[count30minsStats] = current30minsStats
		count30minsStats += 1
	}

	for i := common.Mins30inHrs4; i <= len(min30statsRange)+mins30loopIncrement; i += common.Mins30inHrs4 {
		end := i
		if i > len(min30statsRange) {
			end = len(min30statsRange)
		}
		var statsToCalc common.StatsSet = min30statsRange[i-common.Mins30inHrs4 : end]

		current4hrsStats := statsToCalc.CalcStats()
		hrs4statsRange[count4hrsStats] = current4hrsStats
		count4hrsStats += 1
	}

	for i := common.Mins5inMins30hrs4inHrs24; i <= len(hrs4statsRange)+hrs4loopIncrement; i += common.Mins5inMins30hrs4inHrs24 {
		end := i
		if i > len(hrs4statsRange) {
			end = len(hrs4statsRange)
		}
		var statsToCalc common.StatsSet = hrs4statsRange[i-common.Mins5inMins30hrs4inHrs24 : end]

		current4hrsStats := statsToCalc.CalcStats()
		hrs24statsRange[count24hrsStats] = current4hrsStats
		count24hrsStats += 1
	}

	return models.ResultsByTime{
		Mins5:  mins5statsRange,
		Mins30: min30statsRange,
		Hrs4:   hrs4statsRange,
		Hrs24:  hrs24statsRange,
	}
}
