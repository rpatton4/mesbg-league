package main

import (
	"github.com/rpatton4/mesbg-league/participants/internal/controller/participants"
	handlerhttp "github.com/rpatton4/mesbg-league/participants/internal/handler/http"
	"github.com/rpatton4/mesbg-league/participants/internal/repository/memory"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))

	slog.Info("Starting the Participants service...")
	repo := memory.New()
	ctrl := participants.New(repo)
	handler := handlerhttp.New(ctrl)

	mux := http.NewServeMux()
	mux.Handle("/participants/{id}", http.HandlerFunc(handler.DemuxWithID))
	mux.Handle("/participants", http.HandlerFunc(handler.Demux))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}
}
