package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	appID := mustAtoi(os.Getenv("APP_ID"))
	appHash := os.Getenv("APP_HASH")

	if appID == 0 || appHash == "" {
		log.Fatal("set APP_ID and APP_HASH")
	}

	log.Printf("📱 Testing with APP_ID: %d", appID)

	client := telegram.NewClient(appID, appHash, telegram.Options{})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Connecting to Telegram...")

	if err := client.Run(ctx, func(ctx context.Context) error {
		log.Println("Connected!")

		api := client.API()

		// Пробуем простой запрос
		log.Println("📡 Testing API call...")
		cfg, err := api.HelpGetConfig(ctx)
		if err != nil {
			return fmt.Errorf("API call failed: %w", err)
		}

		log.Printf("API works! DC: %d", cfg.ThisDC)
		return nil
	}); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	log.Println("Test passed!")
}

func mustAtoi(s string) int {
	var v int
	fmt.Sscan(s, &v)
	return v
}
