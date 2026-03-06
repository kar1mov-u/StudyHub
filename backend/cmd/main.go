package main

import (
	"StudyHub/internal/auth"
	"StudyHub/internal/aws"
	"StudyHub/internal/config"
	"StudyHub/internal/gemini"
	"StudyHub/internal/http"
	"StudyHub/internal/modules"
	"StudyHub/internal/rabbitmq"
	"StudyHub/internal/resources"
	studycontent "StudyHub/internal/study_content"
	"StudyHub/internal/users"
	"StudyHub/pgk/postgres"
	"context"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()
	log.Println("staring the v1.3  ...")

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

	s3Storage := aws.NewS3Storage(cfg.BucketName, cfg.AWS_S3_URL)
	geminiClient := gemini.NewGeminiClient(cfg.GeminiKey)

	rbmq := rabbitmq.New(cfg.RBMQUser, cfg.RBMQPass, cfg.RBMQHost)

	//createing srvs
	moduleSrv := modules.NewModuleService(moduleRepo, weeksRepo, moduleRunRepo, academicCalRepo)
	userSrv := users.NewUserService(userRepo)
	authSrv := auth.NewAuthSerivce("", userRepo)
	resourceSrv := resources.NewResourceService(resourceRepo, s3Storage, rbmq)
	contentSrv := studycontent.NewStudyContentService(rbmq, s3Storage, geminiClient)

	httpServer := http.NewHTTPServer(moduleSrv, userSrv, authSrv, resourceSrv, contentSrv, ":8080")

	log.Println("listening...")
	httpServer.Start()
}
