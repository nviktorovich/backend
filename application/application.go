package application

import (
	"context"
	"github.com/NViktorovich/cryptobackend/internal/adapters/client"
	"github.com/NViktorovich/cryptobackend/internal/adapters/storage/postgres"
	"github.com/NViktorovich/cryptobackend/internal/cases"
	"github.com/NViktorovich/cryptobackend/internal/port/server"
	"github.com/NViktorovich/cryptobackend/pkg/client/cryptocompare"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func Run() {
	var CryptoCompareClient client.Scouter
	CryptoCompareClient = &cryptocompare.CryptoCompare{}

	var Client cases.Client
	Client, err := client.NewClientService(CryptoCompareClient)
	if err != nil {
		panic(err)
	}
	if err = godotenv.Load(".env"); err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	connCfg := os.Getenv("PG_CONNECT")
	var Postgres cases.Storage
	Postgres, err = postgres.NewPostgresStorage(connCfg)
	if err != nil {
		panic(err)
	}
	var Service server.Service
	Service, err = cases.NewService(Postgres, Client)
	ctx := context.Background()
	var updatingPeriod time.Duration
	updatingPeriod = 300
	go func() {
		ticker := time.NewTicker(updatingPeriod * time.Second)
		for {
			select {
			case <-ticker.C:
				Service.WriteToStorage(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	var Server *server.Server
	Server, err = server.NewServer(&Service)
	if err != nil {
		panic(err)
	}
	Server.Run()
}
