package main

import (
	"github.com/rpatton4/mesbg-league/games/internal/controller/games"
	handlerhttp "github.com/rpatton4/mesbg-league/games/internal/handler/http"
	"github.com/rpatton4/mesbg-league/games/internal/repository/memory"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))

	slog.Info("Starting the Games service...")
	repo := memory.New()
	ctrl := games.New(repo)
	handler := handlerhttp.New(ctrl)

	mux := http.NewServeMux()
	mux.Handle("/games/{id}", http.HandlerFunc(handler.GetByID))
	mux.Handle("/games", http.HandlerFunc(handler.Demux))
	if err := http.ListenAndServe(":8081", mux); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}

}
