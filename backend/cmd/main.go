package main

import (
	"StudyHub/backend/internal/config"
	"StudyHub/backend/internal/http"
	"StudyHub/backend/internal/modules"
	"StudyHub/backend/pgk/postgres"
	"context"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool := postgres.New(ctx, dbConnString)
	//createing repo's
	moduleRepo := modules.NewModuleRepositoryPostgres(pool)
	moduleRunRepo := modules.NewModuleRunRepositoryPostgres(pool)
	weeksRepo := modules.NewWeekRepositoryPostgres(pool)
	academicCal := modules.NewAcademicCalendarRepositoryPostgres(pool)

	//createing srvs
	moduleSrv := modules.NewModuleService(moduleRepo, weeksRepo, moduleRunRepo, academicCal)

	httpServer := http.NewHTTPServer(moduleSrv, ":8080")

	log.Println("listening...")
	httpServer.Start()

}
