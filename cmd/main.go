package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strconv"
	"stress-relief-ai-chat-back/internal/adapters/http"
	"stress-relief-ai-chat-back/internal/adapters/openai"
	"stress-relief-ai-chat-back/internal/adapters/supabase/users"
	"stress-relief-ai-chat-back/internal/adapters/zap"
	"stress-relief-ai-chat-back/internal/app/chat"
	"syscall"
	"time"
)

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found")
		}
	} else {
		log.Println("Running in Railway environment: ", os.Getenv("RAILWAY_ENVIRONMENT"))
	}

	fmt.Println("Port: ", os.Getenv("PORT"))

	logger, err := zap.NewLogger("development")
	if err != nil {
		log.Fatalf("Error initializing logger: %s", err)
	}

	// Initialize adapters
	openaiAdapter := openai.NewOpenAIAdapter(os.Getenv("OPENAI_API_KEY"),
		os.Getenv("OPENAI_ASSISTANT_ID"),
		logger)

	// Create user storage
	userAPIHandler, err := users.NewUserAPIHandler(os.Getenv("SUPABASE_API_KEY"), os.Getenv("SUPABASE_URL"), logger)
	if err != nil {
		logger.Fatal(context.Background(), "could not create user storage", "error", err.Error())
	}

	// Initialize application services
	chatService := chat.NewChatService(openaiAdapter, logger, userAPIHandler)

	// Setup HTTP server
	server := http.New()
	server.Use(cors.New())

	// Initialize HTTP handlers
	httpHandler := http.NewHandler(chatService, logger)
	httpHandler.SetupRoutes(server.App)

	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			logger.Fatal(context.Background(), "could not parse port", "error", err.Error())
		}
		fmt.Println("Listening on port", port)
		err = server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(fiberServer *http.FiberServer) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
}
