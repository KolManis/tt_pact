package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	api "github.com/KolManis/tt_pact/internal/api/pact/telegram/v1"
	"github.com/KolManis/tt_pact/internal/repository/memory"
	telegramService "github.com/KolManis/tt_pact/internal/service/telegram"
	telegramV1 "github.com/KolManis/tt_pact/pkg/proto/pact/telegram/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
)

func main() {
	repo := memory.NewSessionRepo()
	service := telegramService.NewService(repo)

	api := api.NewAPI(service)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	telegramV1.RegisterTelegramServiceServer(s, api)

	reflection.Register(s)
	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
