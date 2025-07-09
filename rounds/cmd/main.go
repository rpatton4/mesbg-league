package main

import (
	"github.com/rpatton4/mesbg-league/rounds/internal/controller/rounds"
	handlerhttp "github.com/rpatton4/mesbg-league/rounds/internal/handler/http"
	"github.com/rpatton4/mesbg-league/rounds/internal/repository/memory"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Rounds service...")
	repo := memory.New()
	ctrl := rounds.New(repo)
	handler := handlerhttp.New(ctrl)

	http.Handle("/players", http.HandlerFunc(handler.GetRound))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
