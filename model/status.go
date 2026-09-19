package model

type EventStatus string

const (
	EventOpen     EventStatus = "aberto"
	EventClosed   EventStatus = "fechado"
	EventFinished EventStatus = "finalizado"
)

func (s EventStatus) Valid() bool {
	return s == EventOpen || s == EventClosed || s == EventFinished
}

type BetStatus string

const (
	BetPending BetStatus = "pendente"
	BetWon     BetStatus = "ganhou"
	BetLost    BetStatus = "perdeu"
)
