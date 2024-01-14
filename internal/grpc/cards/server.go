package cardsgrpc

import (
	"context"

	cardsv1 "github.com/slavkluev/gophkeeper-contracts/gen/go/cards"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cards/internal/domain/models"
)

type Cards interface {
	GetAll(ctx context.Context) (cards []models.Card, err error)
	SaveCard(
		ctx context.Context,
		number string,
		cvv string,
		month string,
		year string,
		info string,
	) (cardID uint64, err error)
	UpdateCard(
		ctx context.Context,
		id uint64,
		number string,
		cvv string,
		month string,
		year string,
		info string,
	) (err error)
}

type serverAPI struct {
	cardsv1.UnimplementedCardsServer
	cards Cards
}

func Register(gRPCServer *grpc.Server, cards Cards) {
	cardsv1.RegisterCardsServer(gRPCServer, &serverAPI{cards: cards})
}

func (s *serverAPI) GetAll(
	ctx context.Context,
	in *cardsv1.GetAllRequest,
) (*cardsv1.GetAllResponse, error) {
	cards, err := s.cards.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get all cards")
	}

	var crds []*cardsv1.Card
	for _, card := range cards {
		crds = append(crds, &cardsv1.Card{
			Id:     card.ID,
			Number: card.Number,
			Cvv:    card.CVV,
			Month:  card.Month,
			Year:   card.Year,
			Info:   card.Info,
		})
	}

	return &cardsv1.GetAllResponse{Cards: crds}, nil
}

func (s *serverAPI) Save(
	ctx context.Context,
	in *cardsv1.SaveRequest,
) (*cardsv1.SaveResponse, error) {
	if in.GetNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "number is required")
	}

	if in.GetCvv() == "" {
		return nil, status.Error(codes.InvalidArgument, "cvv is required")
	}

	if in.GetMonth() == "" {
		return nil, status.Error(codes.InvalidArgument, "month is required")
	}

	if in.GetYear() == "" {
		return nil, status.Error(codes.InvalidArgument, "year is required")
	}

	accountID, err := s.cards.SaveCard(ctx, in.GetNumber(), in.GetCvv(), in.GetMonth(), in.GetYear(), in.GetInfo())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save card")
	}

	return &cardsv1.SaveResponse{Id: accountID}, nil
}

func (s *serverAPI) Update(
	ctx context.Context,
	in *cardsv1.UpdateRequest,
) (*cardsv1.UpdateResponse, error) {
	if in.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if in.GetNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "number is required")
	}

	if in.GetCvv() == "" {
		return nil, status.Error(codes.InvalidArgument, "cvv is required")
	}

	if in.GetMonth() == "" {
		return nil, status.Error(codes.InvalidArgument, "month is required")
	}

	if in.GetYear() == "" {
		return nil, status.Error(codes.InvalidArgument, "year is required")
	}

	err := s.cards.UpdateCard(ctx, in.GetId(), in.GetNumber(), in.GetCvv(), in.GetMonth(), in.GetYear(), in.GetInfo())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update card")
	}

	return &cardsv1.UpdateResponse{}, nil
}
