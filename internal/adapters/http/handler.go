package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"stress-relief-ai-chat-back/internal/domain"
	"stress-relief-ai-chat-back/internal/ports"
	"strings"
)

type Handler struct {
	chatService ports.ChatService
	logger      ports.Logger
	validator   *validator.Validate
}

func NewHandler(chatService ports.ChatService, logger ports.Logger) *Handler {
	h := &Handler{
		chatService: chatService,
		logger:      logger,
		validator:   validator.New(),
	}
	if h.chatService == nil {
		panic("Cannot create handler without a ChatService")
	}
	if h.logger == nil {
		panic("Cannot create handler without a Logger")
	}
	return h
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Chat routes
	chat := api.Group("/messages")
	chat.Use(h.authMiddleware)
	chat.Post("/", h.handleMessage)
}

func (h *Handler) authMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization token")
	}

	token = strings.TrimPrefix(token, "Bearer ")

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		// Return the secret key for validation
		return []byte(os.Getenv("SUPABASE_JWT_SECRET")), nil
	})

	if err != nil || !parsedToken.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	c.Locals("user", claims)
	c.Locals("userID", claims["sub"])
	return c.Next()
}

func (h *Handler) handleMessage(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	chM := &domain.ChatMessage{
		Content: req.Message,
	}

	userID, ok := c.Locals("userID").(string)
	if !ok {
		h.logger.Error(c.Context(), "could not get user UserID from context")
		return fiber.NewError(fiber.StatusInternalServerError, "Oops! Something went wrong")
	}
	resp, err := h.chatService.ProcessMessage(c.Context(), chM, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(resp)
}
