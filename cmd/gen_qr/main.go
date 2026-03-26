package main

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run gen_qr.go <tg_url>")
		return
	}

	url := os.Args[1]
	err := qrcode.WriteFile(url, qrcode.Medium, 512, "qr.png")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("QR code saved to qr.png")
}
