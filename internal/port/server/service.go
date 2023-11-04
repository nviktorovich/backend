package server

import (
	"context"

	"github.com/NViktorovich/cryptobackend/internal/entities"
)

type Service interface {
	GetAll(ctx context.Context) ([]*entities.Crypto, error)
	GetSpecial(ctx context.Context, title string) (*entities.Crypto, error)
	WriteToStorage(ctx context.Context) error
}
