package common

type Point struct {
	Rate      float64
	Timestamp uint64
}

type PointsSet []Point

type Stats struct {
	Average float64
	High    float64
	Low     float64
	Open    float64
	Close   float64
}

type StatsSet []Stats

func (ds PointsSet) CalcPoints() (stats Stats) {
	var total float64

	stats = Stats{
		Average: 0,
		Open:    ds[len(ds)-1].Rate,
		Close:   ds[0].Rate,
		High:    ds[0].Rate,
		Low:     ds[0].Rate,
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

func (ds StatsSet) CalcStats() (stats Stats) {
	var total float64

	stats = Stats{
		Average: 0,
		High:    ds[0].High,
		Low:     ds[0].Low,
		Open:    ds[len(ds)-1].Open,
		Close:   ds[0].Close,
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
