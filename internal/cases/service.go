package cases

import (
	"context"
	"fmt"
	"github.com/NViktorovich/cryptobackend/internal/entities"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strings"
)

type Service struct {
	storage Storage
	client  Client
	logger  *zap.Logger
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
		logger:  lg,
		tracer:  tr,
	}
	return service, nil
}

func (s *Service) WriteToStorage(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "service: write to storage")
	defer span.End()

	list, err := s.storage.GetList(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "get list failed: %v", err)
		s.logger.Error(err.Error())
		return err
	}

	currentRates, err := s.client.GetCurrentRate(ctx, list)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "get current rates failed: %v", err)
		s.logger.Error(err.Error())
		return err
	}

	if err = s.storage.Write(ctx, currentRates); err != nil {
		err = errors.Wrapf(entities.ErrInternal, "write current rates to the storage failed: %v", err)
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]*entities.Crypto, error) {
	ctx, span := s.tracer.Start(ctx, "service: get all known crypto from storage")
	defer span.End()

	cryptos, err := s.storage.GetAll(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "get all cryptos from storage failed: %v", err)
		s.logger.Error(err.Error())
		return nil, err
	}
	return cryptos, nil
}

func (s *Service) GetSpecial(ctx context.Context, title string) (*entities.Crypto, error) {
	ctx, span := s.tracer.Start(ctx, fmt.Sprintf("service: get special crypto by name: %s", title))
	defer span.End()

	titleList, err := s.storage.GetList(ctx)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "ger list of titles failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	if s.isExist(title, titleList) {
		return s.getExistingSpecialCrypto(ctx, title)
	}

	return s.getMissingSpecialCrypto(ctx, title)
}

func (s *Service) getExistingSpecialCrypto(ctx context.Context, title string) (*entities.Crypto, error) {
	ctx, span := s.tracer.Start(ctx, fmt.Sprintf("service: get crypto from storage by name: %s", title))
	defer span.End()

	crypto, err := s.storage.GetByTitle(ctx, title)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "get crypto from storage by name: %s failed: %v", title, err)
		s.logger.Error(err.Error())
		return nil, err
	}
	return crypto, nil
}

func (s *Service) getMissingSpecialCrypto(ctx context.Context, title string) (*entities.Crypto, error) {
	ctx, span := s.tracer.Start(ctx, fmt.Sprintf("service: get crypto from storage by name: %s", title))
	defer span.End()

	crypto, err := s.client.GetSpecialRate(ctx, title)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "get crypto with special title: %s failed: %v", title, err)
		span.RecordError(err)
		return nil, err
	}

	if err = s.storage.Write(ctx, []*entities.Crypto{crypto}); err != nil {
		err = errors.Wrapf(entities.ErrInternal, "wrati special title: %s to storage failed: %v", title, err)
		span.RecordError(err)
		return nil, err
	}
	return crypto, nil
}

func (s *Service) isExist(specialTitle string, titles []string) bool {
	for _, title := range titles {
		if strings.EqualFold(specialTitle, title) {
			return true
		}
	}
	return false
}
