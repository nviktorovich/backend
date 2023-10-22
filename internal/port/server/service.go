package server

import (
	"context"

	"github.com/NViktorovich/cryptobackend/internal/entities"
)

type Service interface {
	GetLastCrypto(ctx context.Context) ([]entities.Crypto, error)
	GetAvgCrypto(ctx context.Context) ([]entities.Crypto, error)
	GetMinCrypto(ctx context.Context) ([]entities.Crypto, error)
	GetMaxCrypto(ctx context.Context) ([]entities.Crypto, error)
}
