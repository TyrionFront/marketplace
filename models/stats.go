package models

type Stats struct {
	Timestamp uint64  `json:"timestamp" required:"true"`
	Average   float64 `json:"average" required:"true"`
	High      float64 `json:"high" required:"true"`
	Low       float64 `json:"low" required:"true"`
	Open      float64 `json:"open" required:"true"`
	Close     float64 `json:"close" required:"true"`
	TimeFrame string  `json:"time_frame,omitempty"`
}

type StoredStatsDB struct {
	Id        int     `json:"id" required:"true"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Timestamp string  `json:"timestamp" required:"true"`
	Average   float64 `json:"average" required:"true"`
	High      float64 `json:"high" required:"true"`
	Low       float64 `json:"low" required:"true"`
	Open      float64 `json:"open" required:"true"`
	Close     float64 `json:"close" required:"true"`
	TimeFrame string  `json:"time_frame" required:"true"`
}
