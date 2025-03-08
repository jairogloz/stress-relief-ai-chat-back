package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"stress-relief-ai-chat-back/internal/adapters/http"
	"stress-relief-ai-chat-back/internal/adapters/openai"
	"stress-relief-ai-chat-back/internal/adapters/supabase"
	"stress-relief-ai-chat-back/internal/adapters/zap"
	"stress-relief-ai-chat-back/internal/app"
	"stress-relief-ai-chat-back/internal/app/chat"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	logger, err := zap.NewLogger("development")
	if err != nil {
		log.Fatalf("Error initializing logger: %s", err)
	}

	// Initialize adapters
	openaiAdapter := openai.NewOpenAIAdapter(os.Getenv("OPENAI_API_KEY"),
		os.Getenv("OPENAI_ASSISTANT_ID"),
		logger)
	authAdapter := supabase.NewAuthAdapter(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
	)

	// Initialize application services
	chatService := chat.NewChatService(openaiAdapter, logger)
	authService := app.NewAuthService(authAdapter)

	// Setup HTTP server
	server := fiber.New()
	server.Use(cors.New())

	// Initialize HTTP handlers
	httpHandler := http.NewHandler(chatService, authService)
	httpHandler.SetupRoutes(server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(server.Listen(":" + port))
}
