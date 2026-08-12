package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/config"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/logger"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/process"
	"github.com/Knoblauchpilze/galactic-sovereign/cmd/galactic-sovereign/internal"
	"github.com/gin-gonic/gin"
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
// @servers.url /v1/galactic-sovereign
// @servers.description Base path for the galactic-sovereign API
func main() {
	log := logger.New(os.Stdout)

	gin.SetMode(gin.ReleaseMode)

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

	swaggerRoutes, err := internal.SwaggerEndpoints(conf.Server)
	if err != nil {
		log.Error("Failed to create swagger routes", slog.Any("error", err))
		os.Exit(1)
	}
	for _, route := range swaggerRoutes {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
		}
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
