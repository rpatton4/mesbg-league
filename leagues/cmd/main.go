package main

import (
	"github.com/rpatton4/mesbg-league/leagues/internal/primary"
	"github.com/rpatton4/mesbg-league/leagues/internal/secondary"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Leagues service...")
	repo := secondary.New()
	ctrl := primary.New(repo)
	handler := primary.NewHandler(ctrl)

	http.Handle("/leagues", http.HandlerFunc(handler.GetLeague))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
