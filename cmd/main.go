package main

import (
	"net/http"

	"go.uber.org/fx"

	"bet-api/controller"
	"bet-api/db"
	"bet-api/repository"
	"bet-api/usecase"
)

func main() {
	app := fx.New(
		fx.Provide(
			db.Load,
			db.NewDB,
			fx.Annotate(
				repository.NewPostgresUserRepository,
				fx.As(new(repository.UserRepository)),
			),
			fx.Annotate(
				repository.NewPostgresEventRepository,
				fx.As(new(repository.EventRepository)),
			),
			fx.Annotate(
				repository.NewPostgresBetRepository,
				fx.As(new(repository.BetRepository)),
			),
			usecase.NewService,
			controller.NewController,
			db.NewServer,
		),
		fx.Invoke(func(*http.Server) {}),
	)

	app.Run()
}
