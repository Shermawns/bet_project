package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"bet-api/model"
)

type PostgresUserRepository struct {
	db *sql.DB
}

type PostgresEventRepository struct {
	db *sql.DB
}

type PostgresBetRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func NewPostgresEventRepository(db *sql.DB) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

func NewPostgresBetRepository(db *sql.DB) *PostgresBetRepository {
	return &PostgresBetRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	const query = `
		INSERT INTO users (id, nome, saldo_centavos)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.BalanceCents)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	const query = `
		SELECT id, nome, saldo_centavos
		FROM users
		WHERE id = $1
	`

	user := model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.BalanceCents,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresEventRepository) Create(ctx context.Context, event model.Event) (*model.Event, error) {
	const query = `
		INSERT INTO events (id, nome, odds_mil, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, event.ID, event.Name, event.Odds, event.Status)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *PostgresEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	const query = `
		SELECT id, nome, odds_mil, status, time_a_venceu
		FROM events
		WHERE id = $1
	`

	return scanEvent(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresEventRepository) Update(ctx context.Context, id uuid.UUID, odds *model.Odds, status *model.EventStatus) (*model.Event, error) {
	const query = `
		UPDATE events
		SET
			odds_mil = COALESCE($2, odds_mil),
			status = COALESCE($3, status)
		WHERE id = $1
		RETURNING id, nome, odds_mil, status, time_a_venceu
	`

	return scanEvent(r.db.QueryRowContext(ctx, query, id, odds, status))
}

func (r *PostgresEventRepository) List(ctx context.Context, status *model.EventStatus) ([]model.Event, error) {
	const query = `
		SELECT id, nome, odds_mil, status, time_a_venceu
		FROM events
		WHERE ($1::text IS NULL OR status = $1)
		ORDER BY nome
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]model.Event, 0)
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	return events, rows.Err()
}

func scanEvent(row interface{ Scan(...any) error }) (*model.Event, error) {
	event := model.Event{}
	err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Odds,
		&event.Status,
		&event.TeamAWon,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *PostgresBetRepository) Create(ctx context.Context, bet model.Bet) (*model.Bet, error) {
	bet.Status = model.BetPending
	bet.CreatedAt = time.Now().UTC()

	const query = `
		INSERT INTO bets (
			id,
			user_id,
			event_id,
			valor_apostado_centavos,
			odd_no_momento_mil,
			escolha_time_a,
			status,
			created_at
		)
		SELECT
			$1,
			$2,
			e.id,
			$4,
			e.odds_mil,
			$5,
			$6,
			$7
		FROM events e
		WHERE e.id = $3
		RETURNING
			id,
			user_id,
			event_id,
			valor_apostado_centavos,
			odd_no_momento_mil,
			escolha_time_a,
			status,
			created_at
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		bet.ID,
		bet.UserID,
		bet.EventID,
		bet.StakeCents,
		bet.SelectionTeamA,
		bet.Status,
		bet.CreatedAt,
	)

	return scanBet(row)
}

func (r *PostgresBetRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Bet, error) {
	const query = `
		SELECT
			id,
			user_id,
			event_id,
			valor_apostado_centavos,
			odd_no_momento_mil,
			escolha_time_a,
			status,
			created_at
		FROM bets
		WHERE id = $1
	`

	return scanBet(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresBetRepository) List(ctx context.Context) ([]model.Bet, error) {
	const query = `
		SELECT
			id,
			user_id,
			event_id,
			valor_apostado_centavos,
			odd_no_momento_mil,
			escolha_time_a,
			status,
			created_at
		FROM bets
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bets := make([]model.Bet, 0)
	for rows.Next() {
		bet, err := scanBet(rows)
		if err != nil {
			return nil, err
		}
		bets = append(bets, *bet)
	}

	return bets, rows.Err()
}

func (r *PostgresBetRepository) Update(ctx context.Context, id uuid.UUID, stake int64, teamA bool) (*model.Bet, error) {
	const query = `
		UPDATE bets
		SET
			valor_apostado_centavos = $2,
			escolha_time_a = $3
		WHERE id = $1
		RETURNING
			id,
			user_id,
			event_id,
			valor_apostado_centavos,
			odd_no_momento_mil,
			escolha_time_a,
			status,
			created_at
	`

	return scanBet(r.db.QueryRowContext(ctx, query, id, stake, teamA))
}

func (r *PostgresBetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM bets WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return model.ErrNotFound
	}

	return nil
}

func scanBet(row interface{ Scan(...any) error }) (*model.Bet, error) {
	bet := model.Bet{}
	err := row.Scan(
		&bet.ID,
		&bet.UserID,
		&bet.EventID,
		&bet.StakeCents,
		&bet.OddsAtTime,
		&bet.SelectionTeamA,
		&bet.Status,
		&bet.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &bet, nil
}
