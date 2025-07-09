package main

import (
	"github.com/rpatton4/mesbg-league/games/internal/controller/games"
	handlerhttp "github.com/rpatton4/mesbg-league/games/internal/handler/http"
	"github.com/rpatton4/mesbg-league/games/internal/repository/memory"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Games service...")
	repo := memory.New()
	ctrl := games.New(repo)
	handler := handlerhttp.New(ctrl)

	http.Handle("/games", http.HandlerFunc(handler.GetGame))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
