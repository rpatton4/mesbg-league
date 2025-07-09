package main

import (
	"github.com/rpatton4/mesbg-league/leagues/internal/controller/leagues"
	handlerhttp "github.com/rpatton4/mesbg-league/leagues/internal/handler/http"
	"github.com/rpatton4/mesbg-league/leagues/internal/repository/memory"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Leagues service...")
	repo := memory.New()
	ctrl := leagues.New(repo)
	handler := handlerhttp.New(ctrl)

	http.Handle("/leagues", http.HandlerFunc(handler.GetLeague))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
