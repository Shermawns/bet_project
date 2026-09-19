package controller

import (
	"net/http"

	"github.com/google/uuid"

	"bet-api/model"
	"bet-api/usecase"
)

type Controller struct {
	service *usecase.Service
}

func NewController(service *usecase.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Routes() http.Handler {
	mux := http.NewServeMux()

	// Endpoint de verificacao da API.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Endpoints de usuarios.
	mux.HandleFunc("POST /users", c.createUser)
	mux.HandleFunc("GET /users/{id}", c.getUser)

	// Endpoints de eventos.
	mux.HandleFunc("POST /events", c.createEvent)
	mux.HandleFunc("GET /events", c.listEvents)
	mux.HandleFunc("PATCH /events/{id}", c.updateEvent)

	// Endpoints de apostas.
	mux.HandleFunc("POST /bets", c.createBet)
	mux.HandleFunc("GET /bets", c.listBets)
	mux.HandleFunc("GET /bets/{id}", c.getBet)
	mux.HandleFunc("PATCH /bets/{id}", c.updateBet)
	mux.HandleFunc("DELETE /bets/{id}", c.deleteBet)

	return mux
}

type userInput struct {
	Name    string `json:"nome"`
	Balance int64  `json:"saldo_inicial_centavos"`
}

type eventInput struct {
	Name   string             `json:"nome"`
	Odds   string             `json:"odd"`
	Status *model.EventStatus `json:"status"`
}

type betInput struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Stake   int64  `json:"valor_apostado_centavos"`
	TeamA   bool   `json:"escolha_time_a"`
}

func (c *Controller) createUser(w http.ResponseWriter, r *http.Request) {
	var input userInput
	if !decodeJSON(w, r, &input) {
		return
	}

	user, err := c.service.CreateUser(r.Context(), input.Name, input.Balance)
	respond(w, user, err, http.StatusCreated)
}

func (c *Controller) getUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	user, err := c.service.User(r.Context(), id)
	respond(w, user, err, http.StatusOK)
}

func (c *Controller) createEvent(w http.ResponseWriter, r *http.Request) {
	var input eventInput
	if !decodeJSON(w, r, &input) {
		return
	}

	odds, err := model.ParseOdds(input.Odds)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, err)
		return
	}

	event, err := c.service.CreateEvent(r.Context(), input.Name, odds)
	respond(w, eventOutput(event), err, http.StatusCreated)
}

func (c *Controller) updateEvent(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var input eventInput
	if !decodeJSON(w, r, &input) {
		return
	}

	var odds *model.Odds
	if input.Odds != "" {
		value, err := model.ParseOdds(input.Odds)
		if err != nil {
			errorJSON(w, http.StatusBadRequest, err)
			return
		}
		odds = &value
	}

	event, err := c.service.UpdateEvent(r.Context(), id, odds, input.Status)
	respond(w, eventOutput(event), err, http.StatusOK)
}

func (c *Controller) listEvents(w http.ResponseWriter, r *http.Request) {
	var status *model.EventStatus
	if value := r.URL.Query().Get("status"); value != "" {
		parsedStatus := model.EventStatus(value)
		status = &parsedStatus
	}

	events, err := c.service.Events(r.Context(), status)
	if err != nil {
		respond(w, nil, err, http.StatusOK)
		return
	}

	output := make([]eventResponse, 0, len(events))
	for i := range events {
		output = append(output, eventOutput(&events[i]))
	}

	writeJSON(w, http.StatusOK, output)
}

func (c *Controller) createBet(w http.ResponseWriter, r *http.Request) {
	var input betInput
	if !decodeJSON(w, r, &input) {
		return
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, model.ErrInvalidInput)
		return
	}

	eventID, err := uuid.Parse(input.EventID)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, model.ErrInvalidInput)
		return
	}

	bet, err := c.service.CreateBet(r.Context(), userID, eventID, input.Stake, input.TeamA)
	respond(w, betOutput(bet), err, http.StatusCreated)
}

func (c *Controller) listBets(w http.ResponseWriter, r *http.Request) {
	bets, err := c.service.Bets(r.Context())
	if err != nil {
		respond(w, nil, err, http.StatusOK)
		return
	}

	output := make([]betResponse, 0, len(bets))
	for i := range bets {
		output = append(output, betOutput(&bets[i]))
	}

	writeJSON(w, http.StatusOK, output)
}

func (c *Controller) getBet(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	bet, err := c.service.Bet(r.Context(), id)
	respond(w, betOutput(bet), err, http.StatusOK)
}

func (c *Controller) updateBet(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var input betInput
	if !decodeJSON(w, r, &input) {
		return
	}

	bet, err := c.service.UpdateBet(r.Context(), id, input.Stake, input.TeamA)
	respond(w, betOutput(bet), err, http.StatusOK)
}

func (c *Controller) deleteBet(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if err := c.service.DeleteBet(r.Context(), id); err != nil {
		respond(w, nil, err, http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
