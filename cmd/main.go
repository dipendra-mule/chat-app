package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dipendra-mule/chat-app/internal/auth"
	"github.com/dipendra-mule/chat-app/internal/hub"
	"github.com/dipendra-mule/chat-app/pkg/websocket"

	"github.com/rs/cors"
)

func main() {
	// Initialize hub
	chatHub := hub.NewHub()
	go chatHub.Run()

	// WebSocket handler
	wsHandler := websocket.NewHandler(chatHub)

	// HTTP routes
	mux := http.NewServeMux()

	// WebSocket endpoint with auth middleware
	mux.Handle("/ws", auth.AuthMiddleware(http.HandlerFunc(wsHandler.ServeWebSocket)))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Room statistics
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := chatHub.GetRoomStats()
		json.NewEncoder(w).Encode(stats)
	})

	// Serve frontend
	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	// Server configuration
	server := &http.Server{
		Addr:         ":8080",
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")

		// Stop the hub gracefully
		chatHub.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	log.Println("Server stopped")
}
