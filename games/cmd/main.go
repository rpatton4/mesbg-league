package main

import (
	"github.com/rpatton4/mesbg-league/games/internal/inbound"
	"github.com/rpatton4/mesbg-league/games/internal/outbound"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))

	slog.Info("Starting the Games service...")
	repo := outbound.NewMemoryRepository()
	ctrl := inbound.NewHTTPHandler(repo)
	handler := inbound.NewHTTPHandler(ctrl)

	mux := http.NewServeMux()
	mux.Handle("/games/{id}", http.HandlerFunc(handler.DemuxWithID))
	mux.Handle("/games", http.HandlerFunc(handler.Demux))
	if err := http.ListenAndServe(":8081", mux); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}

}
