package cases

import (
	"context"

	"github.com/NViktorovich/cryptobackend/internal/entities"
)

//go:generate mockgen -source=./storage.go -destination=./testdata/storage.go --package=testdata
type Storage interface {
	Write(ctx context.Context, cryptos []entities.Crypto) error
	ReadLast(ctx context.Context, titles []string) ([]entities.Crypto, error)
	ReadAvg(ctx context.Context, titles []string) ([]entities.Crypto, error)
	ReadMin(ctx context.Context, titles []string) ([]entities.Crypto, error)
	ReadMax(ctx context.Context, titles []string) ([]entities.Crypto, error)
	UpdateList(ctx context.Context, title string) error
	GetList(ctx context.Context) ([]string, error)
}
