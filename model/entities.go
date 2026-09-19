package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Odds int64

func ParseOdds(s string) (Odds, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ".")
	if len(parts) > 2 || s == "" {
		return 0, ErrInvalidInput
	}
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 1 {
		return 0, ErrInvalidInput
	}
	fraction := ""
	if len(parts) == 2 {
		fraction = parts[1]
	}
	if len(fraction) > 3 {
		return 0, ErrInvalidInput
	}
	for len(fraction) < 3 {
		fraction += "0"
	}
	frac := int64(0)
	if fraction != "" {
		frac, err = strconv.ParseInt(fraction, 10, 64)
		if err != nil {
			return 0, ErrInvalidInput
		}
	}
	return Odds(whole*1000 + frac), nil
}

func (o Odds) String() string { return fmt.Sprintf("%d.%03d", int64(o)/1000, int64(o)%1000) }
func (o Odds) Valid() bool    { return o >= 1000 }

func (o Odds) Payout(stake int64) int64 { return stake * int64(o) / 1000 }

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"nome"`
	BalanceCents int64     `json:"saldo_centavos"`
}

type Event struct {
	ID       uuid.UUID   `json:"id"`
	Name     string      `json:"nome"`
	Odds     Odds        `json:"-"`
	Status   EventStatus `json:"status"`
	TeamAWon *bool       `json:"time_a_venceu,omitempty"`
}

type Bet struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	EventID        uuid.UUID `json:"event_id"`
	StakeCents     int64     `json:"valor_apostado_centavos"`
	OddsAtTime     Odds      `json:"-"`
	SelectionTeamA bool      `json:"escolha_time_a"`
	Status         BetStatus `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
