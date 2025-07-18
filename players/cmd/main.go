package main

import (
	"github.com/rpatton4/mesbg-league/players/internal/controller/players"
	handlerhttp "github.com/rpatton4/mesbg-league/players/internal/handler/http"
	"github.com/rpatton4/mesbg-league/players/internal/repository/memory"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))

	slog.Info("Starting the Players service...")
	repo := memory.New()
	ctrl := players.New(repo)
	handler := handlerhttp.New(ctrl)

	mux := http.NewServeMux()
	mux.Handle("/players/{id}", http.HandlerFunc(handler.DemuxWithID))
	mux.Handle("/players", http.HandlerFunc(handler.Demux))
	if err := http.ListenAndServe(":8084", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
