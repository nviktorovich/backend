package server

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/NViktorovich/cryptobackend/internal/entities"
	"github.com/NViktorovich/cryptobackend/pkg/dto"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const (
	methodGetCrypto = "/cryptos"
	specialCrypto   = "/{crypto}"
)

var (
	ErrServiceNotSet = errors.New("service not set")
)

type Server struct {
	router  *chi.Mux
	service Service
	logger  *zap.Logger
	tracer  trace.Tracer
}

func NewServer(service Service) (*Server, error) {
	if service == nil {
		return nil, errors.Wrap(entities.ErrInternal, "server creation failed: service is nil")
	}

	lg, err := zap.NewProduction()
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "server creation failed: creating logger: %v", err)
		return nil, err
	}

	tr := otel.Tracer("service")

	s := &Server{
		router:  chi.NewRouter(),
		service: service,
		logger:  lg,
		tracer:  tr,
	}
	return s, nil
}

// @title Simple API
// @version 1.0.0
// @description Simple Crypto API for provided access to information about rate of crypto

// @host localhost:8000
// @BasePath /v1/
func (srv *Server) Run() {

	srv.router.Use(middleware.Logger)

	srv.router.Get(methodGetCrypto, srv.GetAll)
	srv.router.Get(methodGetCrypto+specialCrypto, srv.GetSpecial)

	http.ListenAndServe(":8000", srv.router)
}

// @Summary      all cryptos
// @Description  get data about all known cryptos frob db
// @Tags         crypto
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.Crypto
// @Failure      500  {object} dto.ErrorResponse
// @Router       /cryptos [get]
func (srv *Server) GetAll(rw http.ResponseWriter, req *http.Request) {
	ctx, span := srv.tracer.Start(req.Context(), srv.logger.Name())
	defer span.End()

	res, err := srv.service.GetAll(ctx)
	if err != nil {
		span.RecordError(err)
		srv.sendResponse(rw, http.StatusInternalServerError, err)
	}

	dtoList := make([]*dto.Crypto, 0, len(res))
	for _, crypto := range res {
		dtoList = append(dtoList, srv.convertCryptoToDto(crypto))
	}

	srv.makeSuccessGetResponse(rw, dtoList)
}

// @Summary      special crypto
// @Description  get data about special crypto from db
// @Tags         crypto
// @Accept       json
// @Produce      json
// @Param        title path string true "crypto title"
// @Success      200  {object} dto.Crypto
// @Failure 	404 {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /cryptos/{title} [get]
func (srv *Server) GetSpecial(rw http.ResponseWriter, req *http.Request) {
	ctx, span := srv.tracer.Start(req.Context(), srv.logger.Name())
	defer span.End()

	title := chi.URLParam(req, "title")
	if !srv.validateTitle(title) {
		err := errors.Wrapf(entities.ErrBadRequest, "validate title from url failed: %s", title)
		span.RecordError(err)
		srv.makeErrorResponse(rw, http.StatusBadRequest, err)
	}

	res, err := srv.service.GetSpecial(ctx, title)
	if err != nil {
		err = errors.Wrapf(entities.ErrInternal, "get special title of crypto failed: %v", err)
		span.RecordError(err)
		srv.makeErrorResponse(rw, http.StatusInternalServerError, err)
	}

	srv.sendResponse(rw, http.StatusOK, res)

}

func (srv *Server) sendResponse(rw http.ResponseWriter, statusCode int, obj interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)
	if err := json.NewEncoder(rw).Encode(obj); err != nil {
		srv.logger.Error(err.Error())
	}
}

func (srv *Server) handleError(rw http.ResponseWriter, err error) {
	srv.makeErrorResponse(rw, http.StatusInternalServerError, err)
}

func (srv *Server) makeErrorResponse(rw http.ResponseWriter, statusCode int, err error) {
	srv.sendResponse(rw, statusCode,
		dto.ErrorResponse{
			Message: err.Error(),
		},
	)
}

func (srv *Server) makeSuccessGetResponse(rw http.ResponseWriter, data []*dto.Crypto) {
	srv.sendResponse(rw, http.StatusOK, data)
}

func (srv *Server) convertCryptoToDto(e *entities.Crypto) *dto.Crypto {
	return &dto.Crypto{
		Title:      e.Title,
		ShortTitle: e.ShortTitle,
		Cost:       e.Cost,
		Created:    e.Created.Format(time.RFC3339),
	}
}

func (srv *Server) validateTitle(title string) bool {
	switch {
	case strings.TrimSpace(title) == "":
		return false
	case len([]rune(title)) < 3:
		return false
	case len([]rune(title)) > 255:
		return false
	default:
		return true
	}
}
