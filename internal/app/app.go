package app

import (
	"time"

	"go.uber.org/zap"

	grpcapp "cards/internal/app/grpc"
	"cards/internal/service/cards"
	"cards/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *zap.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	secret string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	cardsService := cards.New(log, storage, storage, storage, tokenTTL)

	grpcApp, err := grpcapp.New(log, cardsService, grpcPort, secret)
	if err != nil {
		panic(err)
	}

	return &App{
		GRPCServer: grpcApp,
	}
}
