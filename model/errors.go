package model

import "errors"

var (
	ErrNotFound          = errors.New("recurso nao encontrado")
	ErrInsufficientFunds = errors.New("saldo insuficiente")
	ErrEventNotOpen      = errors.New("evento nao esta aberto")
	ErrEventSettled      = errors.New("evento ja foi finalizado")
	ErrInvalidInput      = errors.New("dados invalidos")
)
