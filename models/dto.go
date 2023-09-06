package models

type ResultsByTime struct {
	Mins5  []Stats `json:"mins5" required:"true"`
	Mins30 []Stats `json:"mins30" required:"true"`
	Hrs4   []Stats `json:"hrs4" required:"true"`
	Hrs24  []Stats `json:"hrs24" required:"true"`
}

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"-"`
}
