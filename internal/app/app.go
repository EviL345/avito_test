package app

import (
	"fmt"
	"github.com/EviL345/avito_test/internal/config"
	"github.com/EviL345/avito_test/internal/database"
	openapi "github.com/EviL345/avito_test/internal/gen"
	pvzv1 "github.com/EviL345/avito_test/internal/grpc/pvz/v1"
	"github.com/EviL345/avito_test/internal/handler"
	"github.com/EviL345/avito_test/internal/metrics"
	"github.com/EviL345/avito_test/internal/repository"
	"github.com/EviL345/avito_test/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	pvzRepo := repository.NewPVZRepository(db)
	receptionRepo := repository.NewReceptionRepository(db)

	userService := service.NewUserService(userRepo)
	pvzService := service.NewPVZService(pvzRepo)
	receptionService := service.NewReceptionService(receptionRepo, db)

	hndlr := handler.New(userService, pvzService, receptionService)
	r := gin.Default()

	openapi.RegisterHandlers(r, hndlr)

	r.Use(metrics.GetMetricsMiddleware())

	serverAddr := fmt.Sprintf("%s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port)
	go func() {
		pvzv1.Start(cfg.GrpcServer.Port, pvzService)
	}()

	go func() {
		metrics.StartMetricsServer(cfg.PrometheusServer.Port)
	}()

	r.Run(serverAddr)
}
