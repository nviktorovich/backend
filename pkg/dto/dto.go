package dto

type Crypto struct {
	Title      string  `json:"title" db:"title"`
	ShortTitle string  `json:"short_title" db:"short_title"`
	Cost       float64 `json:"cost" db:"cost"`
	Created    string  `json:"created" db:"created"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
