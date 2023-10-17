package cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/NViktorovich/cryptobackend/internal/cases"
	"github.com/NViktorovich/cryptobackend/internal/cases/testdata"
	"github.com/NViktorovich/cryptobackend/internal/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	ErrTest = errors.New("test error")
)

func TestUpdateBase_Succes(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}
	cryptos := []entities.Crypto{
		{ShortTitle: "ETH", Cost: 0.1},
		{ShortTitle: "BTC", Cost: 0.2},
	}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().Write(gomock.Any(), cryptos).Return(nil)
	client := testdata.NewMockClient(ctrl)
	client.EXPECT().GetCurrentRate(gomock.Any(), titles).Return(cryptos, nil)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	err = srv.UpdateBase(context.Background())
	require.Nil(t, err)
}

func TestUpdateBase_GetList_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(nil, ErrTest)
	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	err = srv.UpdateBase(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestUpdateBase_GetCurrentRate_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)

	client := testdata.NewMockClient(ctrl)
	client.EXPECT().GetCurrentRate(gomock.Any(), titles).Return(nil, ErrTest)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	err = srv.UpdateBase(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestUpdateBase_Write_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}
	cryptos := []entities.Crypto{
		{ShortTitle: "ETH", Cost: 0.1},
		{ShortTitle: "BTC", Cost: 0.2},
	}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().Write(gomock.Any(), cryptos).Return(ErrTest)

	client := testdata.NewMockClient(ctrl)
	client.EXPECT().GetCurrentRate(gomock.Any(), titles).Return(cryptos, nil)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	err = srv.UpdateBase(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetLastCrypto_Succes(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}
	cryptos := []entities.Crypto{
		{ShortTitle: "ETH", Cost: 0.1},
		{ShortTitle: "BTC", Cost: 0.2},
	}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadLast(gomock.Any(), titles).Return(cryptos, nil)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetLastCrypto(context.Background())
	require.Nil(t, err)
}

func TestGetLastCrypto_Read_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadLast(gomock.Any(), titles).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetLastCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetLastCrypto_GetList_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetLastCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetAvgCrypto_Succes(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}
	cryptos := []entities.Crypto{
		{ShortTitle: "ETH", Cost: 0.1},
		{ShortTitle: "BTC", Cost: 0.2},
	}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadAvg(gomock.Any(), titles).Return(cryptos, nil)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetAvgCrypto(context.Background())
	require.Nil(t, err)
}

func TestGetAvgCrypto_Read_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadAvg(gomock.Any(), titles).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetAvgCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetAvgCrypto_GetList_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetAvgCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetMinCrypto_Succes(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}
	cryptos := []entities.Crypto{
		{ShortTitle: "ETH", Cost: 0.1},
		{ShortTitle: "BTC", Cost: 0.2},
	}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadMin(gomock.Any(), titles).Return(cryptos, nil)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetMinCrypto(context.Background())
	require.Nil(t, err)
}

func TestGetMinCrypto_Read_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadMin(gomock.Any(), titles).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetMinCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetMinCrypto_GetList_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetMinCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetMaxCrypto_Succes(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}
	cryptos := []entities.Crypto{
		{ShortTitle: "ETH", Cost: 0.1},
		{ShortTitle: "BTC", Cost: 0.2},
	}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadMax(gomock.Any(), titles).Return(cryptos, nil)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetMaxCrypto(context.Background())
	require.Nil(t, err)
}

func TestGetMaxCrypto_Read_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	titles := []string{"ETH", "BTC"}

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(titles, nil)
	storage.EXPECT().ReadMax(gomock.Any(), titles).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetMaxCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestGetMaxCrypto_GetList_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().GetList(gomock.Any()).Return(nil, ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	_, err = srv.GetMaxCrypto(context.Background())
	require.ErrorIs(t, err, entities.ErrInternal)
}

func TestUpdateCryptoList_Succes(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := "USDT"

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().UpdateList(gomock.Any(), title).Return(nil)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	err = srv.UpdateCryptoList(context.Background(), title)
	require.Nil(t, err)
}

func TestUpdateCryptoList_UpdateList_Err(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	title := "USDT"

	storage := testdata.NewMockStorage(ctrl)
	storage.EXPECT().UpdateList(gomock.Any(), title).Return(ErrTest)

	client := testdata.NewMockClient(ctrl)

	srv, err := cases.NewService(storage, client)
	require.Nil(t, err)

	err = srv.UpdateCryptoList(context.Background(), title)
	require.ErrorIs(t, err, entities.ErrInternal)
}
