package cases

import (
	"context"

	"github.com/NViktorovich/cryptobackend/internal/entities"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	storage Storage
	client  Client
	logger  zap.Logger
	tracer  trace.Tracer
}

func NewService(s Storage, c Client) (*Service, error) {
	var err error
	if s == nil {
		err = errors.Wrapf(entities.ErrInvalidParam, "make new service failed, storage is: %v", s)
		return nil, err
	}

	if c == nil {
		err = errors.Wrapf(entities.ErrInvalidParam, "make new service failed, client is: %v", c)
		return nil, err
	}

	lg, err := zap.NewProduction()
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "make new service failed, logger: %v", err)
		return nil, err
	}

	tr := otel.Tracer("service")

	service := &Service{
		storage: s,
		client:  c,
		logger:  *lg,
		tracer:  tr,
	}
	return service, nil
}

func (srv *Service) UpdateBase(ctx context.Context) error {
	ctx, span := srv.tracer.Start(ctx, srv.logger.Name())
	defer span.End()

	titles, err := srv.storage.GetList(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "getting name-list failed: %v", err)
		span.RecordError(err)
		return err
	}

	cryptos, err := srv.client.GetCurrentRate(ctx, titles)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "getting current rates failed: %v", err)
		span.RecordError(err)
		return err
	}

	if err = srv.storage.Write(ctx, cryptos); err != nil {
		err = errors.Wrapf(entities.ErrInternal, "saving current rates to storage failed: %v", err)
		span.RecordError(err)
		return err
	}

	return nil
}

func (srv *Service) GetLastCrypto(ctx context.Context) ([]entities.Crypto, error) {
	ctx, span := srv.tracer.Start(ctx, srv.logger.Name())
	defer span.End()

	titles, err := srv.storage.GetList(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "getting name-list failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	res, err := srv.storage.ReadLast(ctx, titles)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "read failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	return res, nil
}

func (srv *Service) GetAvgCrypto(ctx context.Context) ([]entities.Crypto, error) {
	ctx, span := srv.tracer.Start(ctx, srv.logger.Name())
	defer span.End()

	titles, err := srv.storage.GetList(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "getting name-list failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	res, err := srv.storage.ReadAvg(ctx, titles)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "read failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	return res, nil
}

func (srv *Service) GetMinCrypto(ctx context.Context) ([]entities.Crypto, error) {
	ctx, span := srv.tracer.Start(ctx, srv.logger.Name())
	defer span.End()

	titles, err := srv.storage.GetList(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "getting name-list failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	res, err := srv.storage.ReadMin(ctx, titles)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "read failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	return res, nil
}

func (srv *Service) GetMaxCrypto(ctx context.Context) ([]entities.Crypto, error) {
	ctx, span := srv.tracer.Start(ctx, srv.logger.Name())
	defer span.End()

	titles, err := srv.storage.GetList(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "getting name-list failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	res, err := srv.storage.ReadMax(ctx, titles)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "read failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	return res, nil
}

func (srv *Service) UpdateCryptoList(ctx context.Context, title string) error {
	ctx, span := srv.tracer.Start(ctx, srv.logger.Name())
	defer span.End()

	if err := srv.storage.UpdateList(ctx, title); err != nil {
		err = errors.Wrapf(entities.ErrInternal, "update list of crypto titles failed: %v", err)
		span.RecordError(err)
		return err
	}
	return nil
}
