package dto

type Crypto struct {
	Title      string  `json: "title"`
	ShortTitle string  `json: "short_title"`
	Cost       float64 `json: "cost"`
	TimeStamp  string  `json: "timestamp"`
}

type ErrorResponse struct {
	Error Error
}

type Error struct {
	Message string
}
