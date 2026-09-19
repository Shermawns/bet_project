package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"bet-api/model"
	"bet-api/repository"
)

type Service struct {
	users  repository.UserRepository
	events repository.EventRepository
	bets   repository.BetRepository
}

func NewService(users repository.UserRepository, events repository.EventRepository, bets repository.BetRepository) *Service {
	return &Service{users, events, bets}
}

func (s *Service) CreateUser(ctx context.Context, name string, balance int64) (*model.User, error) {
	if strings.TrimSpace(name) == "" || balance < 0 {
		return nil, model.ErrInvalidInput
	}
	return s.users.Create(ctx, model.User{ID: uuid.New(), Name: name, BalanceCents: balance})
}

func (s *Service) User(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *Service) CreateEvent(ctx context.Context, name string, odds model.Odds) (*model.Event, error) {
	if strings.TrimSpace(name) == "" || !odds.Valid() {
		return nil, model.ErrInvalidInput
	}
	return s.events.Create(ctx, model.Event{ID: uuid.New(), Name: name, Odds: odds, Status: model.EventOpen})
}

func (s *Service) UpdateEvent(ctx context.Context, id uuid.UUID, odds *model.Odds, status *model.EventStatus) (*model.Event, error) {
	if odds != nil && !odds.Valid() {
		return nil, model.ErrInvalidInput
	}
	if status != nil && (!status.Valid() || *status == model.EventFinished) {
		return nil, model.ErrInvalidInput
	}
	return s.events.Update(ctx, id, odds, status)
}

func (s *Service) Events(ctx context.Context, status *model.EventStatus) ([]model.Event, error) {
	if status != nil && !status.Valid() {
		return nil, model.ErrInvalidInput
	}
	return s.events.List(ctx, status)
}

func (s *Service) CreateBet(ctx context.Context, user, event uuid.UUID, stake int64, teamA bool) (*model.Bet, error) {
	if stake <= 0 {
		return nil, model.ErrInvalidInput
	}
	return s.bets.Create(ctx, model.Bet{ID: uuid.New(), UserID: user, EventID: event, StakeCents: stake, SelectionTeamA: teamA, Status: model.BetPending})
}

func (s *Service) Bet(ctx context.Context, id uuid.UUID) (*model.Bet, error) {
	return s.bets.GetByID(ctx, id)
}

func (s *Service) Bets(ctx context.Context) ([]model.Bet, error) {
	return s.bets.List(ctx)
}

func (s *Service) UpdateBet(ctx context.Context, id uuid.UUID, stake int64, teamA bool) (*model.Bet, error) {
	if stake <= 0 {
		return nil, model.ErrInvalidInput
	}
	return s.bets.Update(ctx, id, stake, teamA)
}

func (s *Service) DeleteBet(ctx context.Context, id uuid.UUID) error {
	return s.bets.Delete(ctx, id)
}
