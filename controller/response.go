package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"bet-api/model"
)

type eventResponse struct {
	ID       uuid.UUID         `json:"id"`
	Name     string            `json:"nome"`
	Odds     string            `json:"odd"`
	Status   model.EventStatus `json:"status"`
	TeamAWon *bool             `json:"time_a_venceu,omitempty"`
}

type betResponse struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	EventID   uuid.UUID       `json:"event_id"`
	Stake     int64           `json:"valor_apostado_centavos"`
	Odds      string          `json:"odd_no_momento"`
	Selection bool            `json:"escolha_time_a"`
	Status    model.BetStatus `json:"status"`
	CreatedAt string          `json:"created_at"`
}

func pathID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, model.ErrInvalidInput)
		return uuid.Nil, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, output any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(output); err != nil {
		errorJSON(w, http.StatusBadRequest, model.ErrInvalidInput)
		return false
	}
	return true
}

func respond(w http.ResponseWriter, output any, err error, successStatus int) {
	if err == nil {
		writeJSON(w, successStatus, output)
		return
	}

	status := http.StatusInternalServerError
	if errors.Is(err, model.ErrNotFound) {
		status = http.StatusNotFound
	} else if errors.Is(err, model.ErrInvalidInput) {
		status = http.StatusBadRequest
	} else if errors.Is(err, model.ErrInsufficientFunds) ||
		errors.Is(err, model.ErrEventNotOpen) ||
		errors.Is(err, model.ErrEventSettled) {
		status = http.StatusConflict
	}

	errorJSON(w, status, err)
}

func errorJSON(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"erro": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, output any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(output)
}

func eventOutput(event *model.Event) eventResponse {
	if event == nil {
		return eventResponse{}
	}
	return eventResponse{
		ID:       event.ID,
		Name:     event.Name,
		Odds:     event.Odds.String(),
		Status:   event.Status,
		TeamAWon: event.TeamAWon,
	}
}

func betOutput(bet *model.Bet) betResponse {
	if bet == nil {
		return betResponse{}
	}
	return betResponse{
		ID:        bet.ID,
		UserID:    bet.UserID,
		EventID:   bet.EventID,
		Stake:     bet.StakeCents,
		Odds:      bet.OddsAtTime.String(),
		Selection: bet.SelectionTeamA,
		Status:    bet.Status,
		CreatedAt: bet.CreatedAt.Format(time.RFC3339),
	}
}
