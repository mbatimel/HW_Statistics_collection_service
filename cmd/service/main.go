package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbatimel/HW_Statistics_collection_service/internal/config"
	"github.com/mbatimel/HW_Statistics_collection_service/internal/server"
	"gopkg.in/yaml.v2"
)
func loadConfig(path string) (config.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return config.Config{}, err
	}
	defer file.Close()

	var cfg config.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func main() {
	// Load configuration
	cfg, err := loadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	srv, err := server.NewServerConfig(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := srv.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server run failed: %v", err)
		}
	}()

	log.Println("Server started")
	// Wait for interrupt signal to gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := srv.Close(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	select {
	case <-ctxShutDown.Done():
		if ctxShutDown.Err() == context.DeadlineExceeded {
			log.Println("Shutdown timed out")
		}
	default:
		log.Println("Server exited properly")
	}

}