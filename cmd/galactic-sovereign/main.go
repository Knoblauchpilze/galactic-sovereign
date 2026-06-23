package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/config"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/logger"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/process"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	_ "github.com/Knoblauchpilze/galactic-sovereign/api"
	"github.com/Knoblauchpilze/galactic-sovereign/cmd/galactic-sovereign/internal"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "galactic-sovereign-prod.yml"
	}

	return os.Args[1]
}

// @title			Galactic Sovereign API
// @version		1.0
// @description	REST API for the Galactic Sovereign backend service.
// @servers.url /v1
// @servers.description Base path for the galactic-sovereign API
func main() {
	log := logger.New(os.Stdout)

	conf, err := config.Load(determineConfigName(), internal.DefaultConfig())
	if err != nil {
		log.Error("Failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	conn, err := db.New(context.Background(), conf.Database)
	if err != nil {
		log.Error("Failed to create db connection", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	s := internal.CreateGameServer(conf.Server, conn, log)

	swaggerUi := rest.NewRawRoute(http.MethodGet, "/swagger/*", echoSwagger.WrapHandlerV3)
	if err := s.AddRoute(swaggerUi); err != nil {
		log.Error("Failed to register route", slog.String("route", swaggerUi.Path()), slog.Any("error", err))
		os.Exit(1)
	}

	wait, err := process.StartWithSignalHandler(context.Background(), s)
	if err != nil {
		log.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}

	err = wait()
	if err != nil {
		log.Error("Error while serving", slog.Any("error", err))
		os.Exit(1)
	}
}
