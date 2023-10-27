package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"

	"github.com/NViktorovich/cryptobackend/pkg/dto"

	"github.com/NViktorovich/cryptobackend/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PGStorage struct {
	db     *pgxpool.Pool
	logger *zap.Logger
	tracer trace.Tracer
}

func NewPostgresStorage(cfg string) (*PGStorage, error) {
	pool, err := pgxpool.New(context.Background(), cfg)
	defer pool.Close()

	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "creating pgx pool failed: %v", err)
	}
	lg, err := zap.NewProduction()
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "server creation failed: creating logger: %v", err)
		return nil, err
	}

	tr := otel.Tracer("storage")

	return &PGStorage{
		db:     pool,
		logger: lg,
		tracer: tr,
	}, nil
}

func (s *PGStorage) Write(ctx context.Context, cryptos []*entities.Crypto) error {
	ctx, span := s.tracer.Start(ctx, "pg adapter")
	defer span.End()
	errList := make([]string, 0)
	query := `INSERT INTO crypto_box (short_title, cost) VALUES ($1, $2)`
	for _, crypto := range cryptos {
		dto := s.FromCryptoToDto(crypto)
		parameters := []interface{}{dto.ShortTitle, dto.Cost}
		if err := s.WriteRow(ctx, query, parameters); err != nil {
			errList = append(errList, err.Error())
		}
	}

	if len(errList) > 0 {
		err := errors.Wrap(entities.ErrInternal, strings.Join(errList, ", "))
		span.RecordError(err)
		return err
	}

	return nil
}

func (s *PGStorage) GetAll(ctx context.Context) ([]*entities.Crypto, error) {
	ctx, span := s.tracer.Start(ctx, "pg adapter")
	defer span.End()

	query := `SELECT short_title, cost, created FROM crypto_box 
                                               WHERE created in (SELECT MAX(created) FROM crypto_box 
                                                                                     GROUP BY short_title)`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "get all crypto failed: %v", err)
		span.RecordError(err)
		return nil, err
	}
	dtoList := make([]*dto.Crypto, 0)
	for rows.Next() {
		var dto dto.Crypto
		if err = rows.Scan(&dto); err != nil {
			err = errors.Wrapf(entities.ErrInternal, "scaning failed: %v", err)
			span.RecordError(err)
			return nil, err
		}
		dtoList = append(dtoList, &dto)
	}
	errList := make([]string, 0)
	cryptoList := make([]*entities.Crypto, 0)
	for _, dto := range dtoList {
		crypto, err := s.FromDtoToCrypto(dto)
		if err != nil {
			errList = append(errList, err.Error())
		}
		cryptoList = append(cryptoList, crypto)
	}
	if len(errList) > 0 {
		err = errors.Wrapf(entities.ErrInternal,
			"convert from dto to crypto failed: %s", strings.Join(errList, ", "))
	}

	return cryptoList, err

}

func (s *PGStorage) GetByTitle(ctx context.Context, title string) (*entities.Crypto, error) {
	ctx, span := s.tracer.Start(ctx, "pg adapter")
	defer span.End()

	parameters := []interface{}{title}
	query := `SELECT title, short_title, cost, created FROM 
                                             crypto_box WHERE short_title = $1 AND 
                                            created in (SELECT max(created) FROM crypto_box GROUP BY short_title)`
	var dto = new(dto.Crypto)
	row := s.db.QueryRow(ctx, query, parameters...)
	err := row.Scan(&dto)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errors.Wrapf(entities.ErrNotFound, "search by title: %s has not result", title)
			span.RecordError(err)
			return nil, err
		}
		err = errors.Wrapf(entities.ErrInternal, "search by title: %s failed: %v", title, err)
		span.RecordError(err)
		return nil, err
	}
	crypto, err := s.FromDtoToCrypto(dto)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "convert from dto to crypto failed: %v", err)
		span.RecordError(err)
		return nil, err
	}

	return crypto, nil

}

func (s *PGStorage) UpdateList(ctx context.Context, title string) error {
	ctx, span := s.tracer.Start(ctx, "pg adapter")
	defer span.End()

	parameters := []interface{}{title, title}
	selectQuery := `select id from crypto_box where short_title = $1 or title = $2`
	rows, err := s.db.Query(ctx, selectQuery, parameters...)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "check storage by title or short_title: %s failed: %v", title, err)
		span.RecordError(err)
		return err
	}
	cnt := 0
	for rows.Next() {
		cnt++
	}
	if cnt > 0 {
		err = errors.Wrapf(entities.ErrAlreadyExist, "data with title or short_title: %s already exist", title)
		span.RecordError(err)
		return err
	}

	insertQuery := `INSERT INTO crypto_box (title, short_title) VALUES ($1, $2)`
	err = s.WriteRow(ctx, insertQuery, parameters)
	if err != nil {
		err = errors.Wrapf(entities.ErrAlreadyExist, "writing row to table failed: %v", err)
		span.RecordError(err)
		return err
	}
	return nil
}

func (s *PGStorage) GetList(ctx context.Context) ([]string, error) {
	ctx, span := s.tracer.Start(ctx, "pg adapter")
	defer span.End()

	query := `SELECT DISTINCT short_title from crypto_box`
	rows, err := s.db.Query(ctx, query)
	titles := make([]string, 0)

	for rows.Next() {
		var title string
		err = rows.Scan(&titles)
		if err != nil {
			err = errors.Wrap(entities.ErrInternal, "scanning failed")
			span.RecordError(err)
			return nil, err
		}
		titles = append(titles, title)
	}
	return titles, nil
}

func (s *PGStorage) FromCryptoToDto(crypto *entities.Crypto) *dto.Crypto {
	return &dto.Crypto{
		Title:      crypto.Title,
		ShortTitle: crypto.ShortTitle,
		Cost:       crypto.Cost,
	}
}

func (s *PGStorage) FromDtoToCrypto(dto *dto.Crypto) (*entities.Crypto, error) {
	t, err := time.Parse(time.RFC3339, dto.Created)
	if err != nil {
		return nil, err
	}
	return &entities.Crypto{
		Title:      dto.Title,
		ShortTitle: dto.ShortTitle,
		Cost:       dto.Cost,
		Created:    t,
	}, nil
}

func (s *PGStorage) WriteRow(ctx context.Context, query string, parameters []interface{}) error {
	ctx, span := s.tracer.Start(ctx, "pg adapter")
	defer span.End()

	tag, err := s.db.Exec(ctx, query, parameters...)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if tag.RowsAffected() == 0 {
		err = errors.WithStack(entities.ErrAlreadyExist)
		span.RecordError(err)
		return err
	}

	return nil
}
