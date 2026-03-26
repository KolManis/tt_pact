package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	api "github.com/KolManis/tt_pact/internal/api/pact/telegram/v1"
	"github.com/KolManis/tt_pact/internal/repository/memory"
	telegramService "github.com/KolManis/tt_pact/internal/service/telegram"
	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Читаем APP_ID и APP_HASH
	appIDStr := os.Getenv("APP_ID")
	if appIDStr == "" {
		log.Fatal("APP_ID not set in .env")
	}

	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		log.Fatalf("Invalid APP_ID: %v", err)
	}

	appHash := os.Getenv("APP_HASH")
	if appHash == "" {
		log.Fatal("APP_HASH not set in .env")
	}

	log.Printf("Starting with APP_ID: %d", appID)

	// Создаем репозиторий и сервис
	repo := memory.NewSessionRepo()
	service := telegramService.NewService(repo, appID, appHash) // передаем appID и appHash

	api := api.NewAPI(service)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Создаем gRPC сервер
	s := grpc.NewServer()
	telegramV1.RegisterTelegramServiceServer(s, api)
	reflection.Register(s)

	// Запускаем сервер
	go func() {
		log.Printf("gRPC server listening on port %d", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")

	// Останавливаем все сессии
	service.Shutdown(context.Background())

	s.GracefulStop()
	log.Println("Server stopped")
}
