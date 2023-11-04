package cases

import (
	"context"

	"github.com/NViktorovich/cryptobackend/internal/entities"
)

//go:generate mockgen -source=./client.go -destination=./testdata/client.go --package=testdata
type Client interface {
	GetCurrentRate(ctx context.Context, titles []string) ([]*entities.Crypto, error)
}
