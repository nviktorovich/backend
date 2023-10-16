package entities

import (
	"time"

	"github.com/pkg/errors"
)

// Crypto the main entity, containes info and cost by timestamp
type Crypto struct {
	Title      string
	ShortTitle string
	Cost       float64
	TimeStamp  time.Time
}

func NewCrypto(shortTitle string, cost float64) (*Crypto, error) {
	if cost < 0 {
		err := errors.Wrapf(ErrInvalidParam, "new crypro failed with cost: %.2f", cost)
		return nil, err
	}
	crypto := &Crypto{
		Title:      "",
		ShortTitle: shortTitle,
		Cost:       cost,
	}
	return crypto, nil
}

func (crypto *Crypto) SetTitle(title string) {
	crypto.Title = title
}

func (crypto *Crypto) SetTimeStamp(timestamp time.Time) {
	crypto.TimeStamp = timestamp
}
