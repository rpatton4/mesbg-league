package main

import (
	"github.com/rpatton4/mesbg-league/leagues/internal/controller/league"
	handlerhttp "github.com/rpatton4/mesbg-league/leagues/internal/handler/http"
	"github.com/rpatton4/mesbg-league/leagues/internal/repository/memory"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Leagues service...")
	repo := memory.New()
	ctrl := league.New(repo)
	handler := handlerhttp.New(ctrl)

	http.Handle("/players", http.HandlerFunc(handler.GetLeague))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
