package main

import (
	"StudyHub/backend/internal/auth"
	"StudyHub/backend/internal/config"
	"StudyHub/backend/internal/http"
	"StudyHub/backend/internal/modules"
	"StudyHub/backend/internal/resources"
	"StudyHub/backend/internal/users"
	"StudyHub/backend/pgk/postgres"
	"context"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()
	log.Println("staring...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool := postgres.New(ctx, dbConnString)
	//createing repo's
	moduleRepo := modules.NewModuleRepositoryPostgres(pool)
	moduleRunRepo := modules.NewModuleRunRepositoryPostgres(pool)
	weeksRepo := modules.NewWeekRepositoryPostgres(pool)
	academicCalRepo := modules.NewAcademicCalendarRepositoryPostgres(pool)
	userRepo := users.NewUserRepositoryPostgres(pool)
	resourceRepo := resources.NewResourceRepositoryPostgres(pool)

	s3Storage := resources.NewS3Storage(cfg.BucketName, cfg.AWS_S3_URL)

	//createing srvs
	moduleSrv := modules.NewModuleService(moduleRepo, weeksRepo, moduleRunRepo, academicCalRepo)
	userSrv := users.NewUserService(userRepo)
	authSrv := auth.NewAuthSerivce("", userRepo)
	resourceSrv := resources.NewResourceService(resourceRepo, s3Storage)

	httpServer := http.NewHTTPServer(moduleSrv, userSrv, authSrv, resourceSrv, ":8080")

	log.Println("listening...")
	httpServer.Start()
}
