package main

import (
	"github.com/rpatton4/mesbg-league/players/internal/controller/players"
	handlerhttp "github.com/rpatton4/mesbg-league/players/internal/handler/http"
	"github.com/rpatton4/mesbg-league/players/internal/repository/memory"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Players service...")
	repo := memory.New()
	ctrl := players.New(repo)
	handler := handlerhttp.New(ctrl)

	http.Handle("/players", http.HandlerFunc(handler.GetPlayer))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
