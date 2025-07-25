package main

import (
	//"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	padapters "github.com/rpatton4/mesbg-league/games/internal/gamesprimaryadapters"
	pports "github.com/rpatton4/mesbg-league/games/internal/gamesprimaryports"
	sadapters "github.com/rpatton4/mesbg-league/games/internal/gamessecondaryadapters"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))

	slog.Info("Starting the Games service...")
	repo := sadapters.NewMemoryRepository()
	ctrl := pports.NewTxnController(repo)
	handler := padapters.NewHTTPHandler(ctrl)

	router := http.NewServeMux()

	api := humago.New(router, huma.DefaultConfig("Games Service", "1.0.0"))

	huma.Get(api, "/games/{id}", handler.HumaGetByID)
	//mux.Handle("/games/{id}", http.HandlerFunc(handler.DemuxWithID))
	//mux.Handle("/games", http.HandlerFunc(handler.Demux))

	if err := http.ListenAndServe(":8081", router); err != nil {
		slog.Error("Failed to start HTTP server", "error", err.Error())
		panic(err)
	}

}
