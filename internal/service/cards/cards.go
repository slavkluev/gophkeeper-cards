package cards

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"cards/internal/domain/models"
)

type Cards struct {
	log          *zap.Logger
	cardSaver    CardSaver
	cardUpdater  CardUpdater
	cardProvider CardProvider
	tokenTTL     time.Duration
}

type CardSaver interface {
	SaveCard(
		ctx context.Context,
		number string,
		cvv string,
		month string,
		year string,
		info string,
		userUID uint64,
	) (uid uint64, err error)
}

type CardUpdater interface {
	UpdateCard(
		ctx context.Context,
		id uint64,
		number string,
		cvv string,
		month string,
		year string,
		info string,
		userUID uint64,
	) error
}

type CardProvider interface {
	GetAll(ctx context.Context, userUID uint64) ([]models.Card, error)
}

func New(
	log *zap.Logger,
	cardSaver CardSaver,
	cardUpdater CardUpdater,
	cardProvider CardProvider,
	tokenTTL time.Duration,
) *Cards {
	return &Cards{
		cardSaver:    cardSaver,
		cardUpdater:  cardUpdater,
		cardProvider: cardProvider,
		log:          log,
		tokenTTL:     tokenTTL,
	}
}

func (a *Cards) SaveCard(ctx context.Context, number string, cvv string, month string, year string, info string) (uint64, error) {
	const op = "Cards.SaveCard"

	log := a.log.With(
		zap.String("op", op),
		zap.String("number", number),
	)

	log.Info("attempting to save card")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return 0, fmt.Errorf("%s: failed to find user uid", op)
	}

	id, err := a.cardSaver.SaveCard(ctx, number, cvv, month, year, info, userUID)
	if err != nil {
		log.Error("failed to save card", zap.Error(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("card was saved successfully")

	return id, nil
}

func (a *Cards) UpdateCard(ctx context.Context, id uint64, number string, cvv string, month string, year string, info string) error {
	const op = "Cards.UpdateCard"

	log := a.log.With(
		zap.String("op", op),
		zap.Uint64("id", id),
		zap.String("number", number),
	)

	log.Info("attempting to update card")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return fmt.Errorf("%s: failed to find user uid", op)
	}

	err := a.cardUpdater.UpdateCard(ctx, id, number, cvv, month, year, info, userUID)
	if err != nil {
		log.Error("failed to update card", zap.Error(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("card was updated successfully")

	return nil
}

func (a *Cards) GetAll(ctx context.Context) ([]models.Card, error) {
	const op = "Cards.GetAll"

	log := a.log.With(
		zap.String("op", op),
	)

	log.Info("attempting to get all cards")

	rawUserUID := ctx.Value("user-uid")
	userUID, ok := rawUserUID.(uint64)
	if !ok {
		log.Error("failed to find user uid")

		return nil, fmt.Errorf("%s: failed to find user uid", op)
	}

	cards, err := a.cardProvider.GetAll(ctx, userUID)
	if err != nil {
		a.log.Error("failed to get all cards", zap.Error(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("cards are got successfully")

	return cards, nil
}
