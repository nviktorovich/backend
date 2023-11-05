package cases_test

import (
	"context"
	"errors"
	"github.com/NViktorovich/cryptobackend/internal/cases"
	"github.com/NViktorovich/cryptobackend/internal/cases/testdata"
	"github.com/NViktorovich/cryptobackend/internal/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strings"
	"testing"
)

var (
	errTest = errors.New("test error")
)

func TestNewService_NilStorage_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := testdata.NewMockClient(ctrl)
	service, err := cases.NewService(nil, client)
	require.ErrorIs(t, err, entities.ErrInvalidParam)
	require.Nil(t, service)
}

func TestNewService_NilClient_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	service, err := cases.NewService(storage, nil)
	require.ErrorIs(t, err, entities.ErrInvalidParam)
	require.Nil(t, service)
}

func TestNewService_Successful(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)
}

func TestWriteToStorage_GetList_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return(nil, errTest)
	err = service.WriteToStorage(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestWriteToStorage_GetCurrentRate_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	list := []string{makeString(), makeString()}

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return(list, nil)

	client.EXPECT().GetCurrentRate(gomock.Any(), list).Return(nil, errTest)
	err = service.WriteToStorage(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestWriteToStorage_Write_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	list := []string{makeString(), makeString()}
	currentRates := []*entities.Crypto{&entities.Crypto{ShortTitle: makeString()}, &entities.Crypto{ShortTitle: makeString()}}

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return(list, nil)

	client.EXPECT().GetCurrentRate(gomock.Any(), list).Return(currentRates, nil)

	storage.EXPECT().Write(gomock.Any(), currentRates).Return(errTest)

	err = service.WriteToStorage(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestWriteToStorage_Successful(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	list := []string{makeString(), makeString()}
	currentRates := []*entities.Crypto{&entities.Crypto{ShortTitle: makeString()}, &entities.Crypto{ShortTitle: makeString()}}

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return(list, nil)

	client.EXPECT().GetCurrentRate(gomock.Any(), list).Return(currentRates, nil)

	storage.EXPECT().Write(gomock.Any(), currentRates).Return(nil)

	err = service.WriteToStorage(context.Background())
	require.NoError(t, err)
}

func Test_GetAll_GetAll_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetAll(gomock.Any()).Return(nil, errTest)

	res, err := service.GetAll(context.Background())
	require.Nil(t, res)
	require.ErrorIs(t, err, entities.ErrInternal)
}

func Test_GetAll_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rates := []*entities.Crypto{&entities.Crypto{ShortTitle: makeString()}, &entities.Crypto{ShortTitle: makeString()}}

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetAll(gomock.Any()).Return(rates, nil)

	res, err := service.GetAll(context.Background())
	require.NoError(t, err)
	require.Equal(t, rates, res)
}

func Test_GetSpecial_GetList_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := makeString()

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return(nil, errTest)

	res, err := service.GetSpecial(context.Background(), title)
	require.Nil(t, res)
	require.ErrorIs(t, err, entities.ErrInternal)
}

func Test_GetSpecial_GetByTitle_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := makeString()

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return([]string{title}, nil)
	storage.EXPECT().GetByTitle(gomock.Any(), title).Return(nil, errTest)

	res, err := service.GetSpecial(context.Background(), title)
	require.Nil(t, res)
	require.ErrorIs(t, err, entities.ErrInternal)
}

func Test_GetSpecial_GetByTitle_getExistingSpecialCrypto_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := makeString()

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return([]string{title}, nil)
	storage.EXPECT().GetByTitle(gomock.Any(), title).Return(&entities.Crypto{ShortTitle: title}, nil)

	res, err := service.GetSpecial(context.Background(), title)
	require.NotNil(t, res)
	require.NoError(t, err)
}

func Test_GetSpecial_getMissingSpecialCrypto_GetCurrentRate_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := makeString()

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return([]string{makeString()}, nil)
	client.EXPECT().GetCurrentRate(gomock.Any(), []string{title}).Return(nil, errTest)

	res, err := service.GetSpecial(context.Background(), title)
	require.Nil(t, res)
	require.ErrorIs(t, err, entities.ErrInternal)
}

func Test_GetSpecial_getMissingSpecialCrypto_Write_Err(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := makeString()
	crypto := &entities.Crypto{ShortTitle: title}
	cryptos := []*entities.Crypto{crypto}

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return([]string{makeString()}, nil)
	client.EXPECT().GetCurrentRate(gomock.Any(), []string{title}).Return(cryptos, nil)
	storage.EXPECT().Write(gomock.Any(), cryptos).Return(errTest)

	res, err := service.GetSpecial(context.Background(), title)
	require.Nil(t, res)
	require.ErrorIs(t, err, entities.ErrInternal)
}

func Test_GetSpecial_getMissingSpecialCrypto_Successful(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := makeString()
	crypto := &entities.Crypto{ShortTitle: title}
	cryptos := []*entities.Crypto{crypto}

	storage := testdata.NewMockStorage(ctrl)
	client := testdata.NewMockClient(ctrl)

	service, err := cases.NewService(storage, client)
	require.NoError(t, err)
	require.NotNil(t, service)

	storage.EXPECT().GetList(gomock.Any()).Return([]string{makeString()}, nil)
	client.EXPECT().GetCurrentRate(gomock.Any(), []string{title}).Return(cryptos, nil)
	storage.EXPECT().Write(gomock.Any(), cryptos).Return(nil)

	res, err := service.GetSpecial(context.Background(), title)
	require.NotNil(t, res)
	require.NoError(t, err)
}

func makeString() string {
	var s strings.Builder
	for i := 0; i < 10; i++ {
		symbol := rand.Intn(127)
		s.WriteByte(byte(symbol))
	}
	return s.String()
}
