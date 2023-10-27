package cases

import (
	"context"

	"github.com/NViktorovich/cryptobackend/internal/entities"
)

//go:generate mockgen -source=./storage.go -destination=./testdata/storage.go --package=testdata
type Storage interface {
	Write(ctx context.Context, cryptos []*entities.Crypto) error
	GetAll(ctx context.Context) ([]*entities.Crypto, error)
	GetByTitle(ctx context.Context, title string) (*entities.Crypto, error)
	GetList(ctx context.Context) ([]string, error)
}
