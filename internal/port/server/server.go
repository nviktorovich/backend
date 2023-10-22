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
	ParamService = "service"
	BaseURL      = "/crypto/v1"

	methodGetLast = "getLast"
	methodGetAvg  = "getAvg"
	methodGetMin  = "getMin"
	methodGetMax  = "getMax"
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

func (serv *Server) Run() {

	serv.router.Use(middleware.Logger)

	serv.router.Get(methodGetLast, serv.GetLast)
	serv.router.Get(methodGetAvg, serv.GetLast)
	serv.router.Get(methodGetMin, serv.GetLast)
	serv.router.Get(methodGetMax, serv.GetLast)

	http.ListenAndServe(":8000", serv.router)
}

// @Summary      last rates known cryptos
// @Description  last rates known cryptos
// @Tags         crypto
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.Crypto
// @Failure      500  {object} dto.Error
// @Router       /getLast [get]
func (serv *Server) GetLast(rw http.ResponseWriter, req *http.Request) {
	ctx, span := serv.tracer.Start(req.Context(), serv.logger.Name())
	res, err := serv.service.GetLastCrypto(ctx)
	if err != nil {
		span.RecordError(err)
		serv.sendResponse(rw, http.StatusInternalServerError, err)
	}
	dtoRes := serv.ConvertCryptoToDtoCrypto(res)
	serv.makeSuccessGetResponse(rw, dtoRes)
}

// @Summary      avg rates known cryptos
// @Description  avg rates known cryptos
// @Tags         crypto
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.Crypto
// @Failure      500  {object} dto.Error
// @Router       /getAvg [get]
func (serv *Server) GetAvg(rw http.ResponseWriter, req *http.Request) {
	ctx, span := serv.tracer.Start(req.Context(), serv.logger.Name())
	res, err := serv.service.GetAvgCrypto(ctx)
	if err != nil {
		span.RecordError(err)
		serv.sendResponse(rw, http.StatusInternalServerError, err)
	}
	dtoRes := serv.ConvertCryptoToDtoCrypto(res)
	serv.makeSuccessGetResponse(rw, dtoRes)
}

// @Summary      min rates known cryptos
// @Description  min rates known cryptos
// @Tags         crypto
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.Crypto
// @Failure      500  {object} dto.Error
// @Router       /getMin [get]
func (serv *Server) GetMin(rw http.ResponseWriter, req *http.Request) {
	ctx, span := serv.tracer.Start(req.Context(), serv.logger.Name())
	res, err := serv.service.GetMinCrypto(ctx)
	if err != nil {
		span.RecordError(err)
		serv.sendResponse(rw, http.StatusInternalServerError, err)
	}
	dtoRes := serv.ConvertCryptoToDtoCrypto(res)
	serv.makeSuccessGetResponse(rw, dtoRes)
}

// @Summary      max rates known cryptos
// @Description  max rates known cryptos
// @Tags         crypto
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.Crypto
// @Failure      500  {object} dto.Error
// @Router       /getMax [get]
func (serv *Server) GetMax(rw http.ResponseWriter, req *http.Request) {
	ctx, span := serv.tracer.Start(req.Context(), serv.logger.Name())
	res, err := serv.service.GetMaxCrypto(ctx)
	if err != nil {
		span.RecordError(err)
		serv.sendResponse(rw, http.StatusInternalServerError, err)
	}
	dtoRes := serv.ConvertCryptoToDtoCrypto(res)
	serv.makeSuccessGetResponse(rw, dtoRes)
}

func (serv *Server) sendResponse(rw http.ResponseWriter, statusCode int, obj interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)
	if err := json.NewEncoder(rw).Encode(obj); err != nil {
		serv.logger.Error(err.Error())
	}
}

func (serv *Server) handleError(rw http.ResponseWriter, err error) {
	serv.makeErrorResponse(rw, http.StatusInternalServerError, err)
}

func (serv *Server) makeErrorResponse(rw http.ResponseWriter, statusCode int, err error) {
	serv.sendResponse(rw, statusCode,
		dto.ErrorResponse{
			Error: dto.Error{
				Message: err.Error(),
			},
		},
	)
}

func (serv *Server) makeSuccessGetResponse(rw http.ResponseWriter, data []dto.Crypto) {
	serv.sendResponse(rw, http.StatusOK, data)
}

func (serv *Server) ConvertCryptoToDtoCrypto(in []entities.Crypto) []dto.Crypto {
	dtoCrytoList := make([]dto.Crypto, 0, len(in))
	for _, crypto := range in {
		dtoCrytoList = append(dtoCrytoList, dto.Crypto{
			Title:      crypto.Title,
			ShortTitle: crypto.ShortTitle,
			Cost:       crypto.Cost,
			TimeStamp:  crypto.TimeStamp.String(),
		})
	}
	return dtoCrytoList
}

func (serv *Server) ConvertDtoCryptoToCrypto(in []dto.Crypto) ([]entities.Crypto, error) {
	cryptoList := make([]entities.Crypto, 0, len(in))
	errList := make([]string, 0)
	for _, crypto := range in {
		t, err := time.Parse(time.RFC3339, crypto.TimeStamp)
		if err != nil {
			errList = append(errList, err.Error())
			continue
		}

		cryptoList = append(cryptoList, entities.Crypto{
			Title:      crypto.Title,
			ShortTitle: crypto.ShortTitle,
			Cost:       crypto.Cost,
			TimeStamp:  t,
		})
	}
	if len(errList) > 0 {
		err := errors.Wrapf(entities.ErrInvalidParam, "convert data to time failed: %v", strings.Join(errList, ", "))
		return cryptoList, err
	}
	return cryptoList, nil
}
