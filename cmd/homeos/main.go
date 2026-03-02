// cmd/homeos/main.go
// HomeOS — Self-Organizing AI Agent Architecture
// Entry point chính

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/goldlotus1810/HomeOS/internal/aam"
	"github.com/goldlotus1810/HomeOS/internal/leoai"
)

func main() {
	log.Println("🏠 HomeOS — Starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Khởi động LeoAI
	leo, err := leoai.New("./data/tree", 512)
	if err != nil {
		log.Fatalf("Failed to start LeoAI: %v", err)
	}
	go leo.Run(ctx)
	log.Println("✅ LeoAI: Online")

	// Khởi động AAM
	ui := &consoleUI{}
	aamCore := aam.New(ui)

	// TODO: Đăng ký Chiefs
	// aamCore.RegisterChief(chiefs.NewHomeChief(...))
	// aamCore.RegisterChief(chiefs.NewVisionChief(...))

	log.Println("✅ AAM: Online — System ready")
	log.Println("📡 Waiting for input...")

	// Chạy cho đến khi nhận signal dừng
	go func() {
		if err := aamCore.Run(ctx); err != nil {
			log.Printf("AAM error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 HomeOS: Shutting down gracefully...")
	cancel()
}

// consoleUI là implementation đơn giản cho console
type consoleUI struct{}

func (u *consoleUI) AskUser(_ context.Context, question string) (string, error) {
	log.Printf("❓ AAM asks: %s", question)
	return "", nil
}

func (u *consoleUI) Notify(_ context.Context, msg string) error {
	log.Printf("📢 AAM: %s", msg)
	return nil
}

func (u *consoleUI) AskConfirm(_ context.Context, question string) (bool, error) {
	log.Printf("🔔 AAM confirms: %s", question)
	return true, nil // Auto-confirm in dev mode
}
