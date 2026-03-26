package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/tg"
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/term"
)

func main() {
	appID := 29640364
	appHash := "10a22687a449dc75183ef4d76f7d986a"
	if appID == 0 || appHash == "" {
		log.Fatal("set APP_ID and APP_HASH")
	}

	ctx := context.Background()

	disp := tg.NewUpdateDispatcher()
	client := telegram.NewClient(appID, appHash, telegram.Options{
		UpdateHandler: disp,
		SessionStorage: &telegram.FileSessionStorage{
			Path: "./session.json",
		},
	})

	if err := client.Run(ctx, func(ctx context.Context) error {
		// Если уже есть валидная сессия — ничего делать не нужно.
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return err
		}
		if status.Authorized {
			log.Println("Already authorized")
			return nil
		}

		loggedIn := qrlogin.OnLoginToken(disp)
		_, err = client.QR().Auth(ctx, loggedIn, func(ctx context.Context, token qrlogin.Token) error {
			// 1) Можно открыть URL напрямую:
			log.Printf("Open URL in Telegram: %s", token.URL())

			// 2) Или сохранить нормальный PNG QR:
			if err := qrcode.WriteFile(token.URL(), qrcode.Medium, 512, "login_qr.png"); err != nil {
				log.Printf("write QR png: %v", err)
			} else {
				log.Printf("QR saved to ./login_qr.png (expires: %s)", token.Expires())
			}
			return nil
		})

		// после QR Telegram может потребовать облачный пароль (2FA).
		if err != nil {
			if errors.Is(err, auth.ErrPasswordAuthNeeded) || strings.Contains(err.Error(), "SESSION_PASSWORD_NEEDED") {
				pass, pErr := readPassword("Enter cloud password (2FA): ")
				if pErr != nil {
					return pErr
				}
				if _, pErr = client.Auth().Password(ctx, pass); pErr != nil {
					return pErr
				}
				log.Println("2FA password accepted")
				return nil
			}
			return err
		}

		log.Println("QR auth successful")
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	// Если хотите брать из env:
	// if p := os.Getenv("PASSWORD"); p != "" { return p, nil }

	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(b))
	if s == "" {
		return "", fmt.Errorf("empty password")
	}
	return s, nil
}

func mustAtoi(s string) int {
	var v int
	_, _ = fmt.Fscan(strings.NewReader(s), &v)
	return v
}
