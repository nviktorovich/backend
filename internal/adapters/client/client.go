package client

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strings"

	"github.com/NViktorovich/cryptobackend/internal/entities"
)

const (
	dollar = "USD"
)

type ClientService struct {
	scouter Scouter
	logger  *zap.Logger
	tracer  trace.Tracer
}

func NewClientService(sc Scouter) (*ClientService, error) {
	if sc == nil {
		err := errors.Wrap(entities.ErrInternal, "created client service failed, scouter is nil")
		return nil, err
	}

	lg, err := zap.NewProduction()
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "client service creation failed: creating logger: %v", err)
		return nil, err
	}

	tr := otel.Tracer("service")

	return &ClientService{
		scouter: sc,
		logger:  lg,
		tracer:  tr,
	}, nil
}

func (cs *ClientService) GetCurrentRate(ctx context.Context, titles []string) ([]*entities.Crypto, error) {
	ctx, span := cs.tracer.Start(ctx, "service: write to storage")
	defer span.End()

	res, err := cs.scouter.GetAll(titles, dollar)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "scouter return error: %v", err)
		span.RecordError(err)
		return nil, err
	}
	cryptos := make([]*entities.Crypto, 0)
	errList := make([]string, 0)
	for title, _ := range res {
		crypto, err := cs.convertMapToCrypto(title, res)
		if err != nil {
			errList = append(errList, err.Error())
			continue
		}
		cryptos = append(cryptos, crypto)
	}
	if len(errList) > 0 {
		err = errors.Wrapf(entities.ErrInternal, "failed creat new crypto: %s", strings.Join(errList, ", "))
		span.RecordError(err)
		return nil, err
	}
	return cryptos, nil
}

func (cs *ClientService) convertMapToCrypto(title string, values map[string]float64) (*entities.Crypto, error) {
	cost, ok := values[title]
	if !ok {
		return nil, errors.WithStack(entities.ErrInternal)
	}
	return entities.NewCrypto(title, cost)
}
