package repository

import (
	"context"

	"github.com/google/uuid"

	"bet-api/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type EventRepository interface {
	Create(ctx context.Context, event model.Event) (*model.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	Update(ctx context.Context, id uuid.UUID, odds *model.Odds, status *model.EventStatus) (*model.Event, error)
	List(ctx context.Context, status *model.EventStatus) ([]model.Event, error)
}

type BetRepository interface {
	Create(ctx context.Context, bet model.Bet) (*model.Bet, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Bet, error)
	List(ctx context.Context) ([]model.Bet, error)
	Update(ctx context.Context, id uuid.UUID, stake int64, teamA bool) (*model.Bet, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
