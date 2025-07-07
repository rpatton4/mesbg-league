package main

import (
	"github.com/rpatton4/mesbg-league/player/internal/controller/player"
	handlerhttp "github.com/rpatton4/mesbg-league/player/internal/handler/http"
	"github.com/rpatton4/mesbg-league/player/internal/repository/memory"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Player service...")
	repo := memory.New()
	ctrl := player.New(repo)
	handler := handlerhttp.New(ctrl)

	http.Handle("/players", http.HandlerFunc(handler.GetPlayer))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
