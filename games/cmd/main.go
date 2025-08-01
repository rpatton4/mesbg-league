package main

import (
	//"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	padapters "github.com/rpatton4/mesbg-league/games/internal/primary"
	sadapters "github.com/rpatton4/mesbg-league/games/internal/secondary"
	"log/slog"
	"net/http"
	"os"
)

var port = "8081"

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))

	slog.Info("Starting the Games service on port " + port)
	repo := sadapters.NewMemoryRepository()
	ctrl := padapters.NewTxnController(repo)
	handler := padapters.NewHumaHandler(ctrl)

	router := http.NewServeMux()

	api := humago.New(router, huma.DefaultConfig("Games Service", "1.0.0"))

	huma.Get(api, "/games/{id}", handler.GetByID)
	huma.Post(api, "/games", handler.Post)
	huma.Put(api, "/games/{id}", handler.Put)
	huma.Delete(api, "/games/{id}", handler.Delete)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}

}
