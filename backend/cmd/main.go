package main

import (
	"StudyHub/internal/auth"
	"StudyHub/internal/aws"
	"StudyHub/internal/comments"
	"StudyHub/internal/config"
	"StudyHub/internal/content"
	"StudyHub/internal/gemini"
	"StudyHub/internal/http"
	"StudyHub/internal/modules"
	"StudyHub/internal/rabbitmq"
	"StudyHub/internal/resources"
	"StudyHub/internal/users"
	"StudyHub/pgk/postgres"
	"context"
	"fmt"
	"log"
)

// user reverse proxy ga borad, keyn nginx url ga qarap front yoki back ligni blad, agar /api busa bu back ga ketad
func main() {
	//load configs
	cfg := config.Load()
	log.Println("staring the v1.3  ...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	//connect to the DB
	pool := postgres.New(ctx, dbConnString)
	//createing repo's
	moduleRepo := modules.NewModuleRepositoryPostgres(pool)
	moduleRunRepo := modules.NewModuleRunRepositoryPostgres(pool)
	weeksRepo := modules.NewWeekRepositoryPostgres(pool)
	academicCalRepo := modules.NewAcademicCalendarRepositoryPostgres(pool)
	userRepo := users.NewUserRepositoryPostgres(pool)
	resourceRepo := resources.NewResourceRepositoryPostgres(pool)
	contentRepo := content.NewContentRepositoryPostgres(pool)
	commentRepo := comments.NewCommentRepositoryPostgres(pool)

	//create instances for external services
	s3Storage := aws.NewS3Storage(cfg.BucketName, cfg.AWS_S3_URL)
	geminiClient := gemini.NewGeminiClient(cfg.GeminiKey)
	rbmq := rabbitmq.New(cfg.RBMQUser, cfg.RBMQPass, cfg.RBMQHost)

	//createing srvs
	moduleSrv := modules.NewModuleService(moduleRepo, weeksRepo, moduleRunRepo, academicCalRepo)
	userSrv := users.NewUserService(userRepo)
	authSrv := auth.NewAuthSerivce("", userRepo)
	resourceSrv := resources.NewResourceService(resourceRepo, s3Storage, rbmq)
	contentSrv := content.NewContentService(contentRepo, rbmq, s3Storage, geminiClient)
	commentSrv := comments.NewCommentService(commentRepo)

	httpServer := http.NewHTTPServer(moduleSrv, userSrv, authSrv, resourceSrv, contentSrv, commentSrv, ":8080")

	log.Println("listening...")
	httpServer.Start()
}
