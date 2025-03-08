package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"stress-relief-ai-chat-back/internal/app"
	"stress-relief-ai-chat-back/internal/domain"
	"stress-relief-ai-chat-back/internal/ports"
)

type Handler struct {
	chatService ports.ChatService
	authService *app.AuthService
	validator   *validator.Validate
}

func NewHandler(chatService ports.ChatService, authService *app.AuthService) *Handler {
	return &Handler{
		chatService: chatService,
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Chat routes
	chat := api.Group("/messages")
	//chat.Use(h.authMiddleware)
	chat.Post("/", h.handleMessage)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/validate", h.handleValidateToken)
}

func (h *Handler) authMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization token")
	}

	user, err := h.authService.ValidateToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	c.Locals("user", user)
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

	//user := c.Locals("user").(*domain.User)
	resp, err := h.chatService.ProcessMessage(c.Context(), chM, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(resp)
}

func (h *Handler) handleValidateToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	return c.JSON(user)
}
