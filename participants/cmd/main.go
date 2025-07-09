package main

import (
	"github.com/rpatton4/mesbg-league/participants/internal/controller/participants"
	handlerhttp "github.com/rpatton4/mesbg-league/participants/internal/handler/http"
	"github.com/rpatton4/mesbg-league/participants/internal/repository/memory"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting the Participants service...")
	repo := memory.New()
	ctrl := participants.New(repo)
	handler := handlerhttp.New(ctrl)

	http.Handle("/participants", http.HandlerFunc(handler.GetParticipant))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
