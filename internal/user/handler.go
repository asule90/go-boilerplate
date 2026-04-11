package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sule/go-boilerplate/config"
	"github.com/sule/go-boilerplate/internal/middleware"
	"github.com/sule/go-boilerplate/internal/utils"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for the user domain.
type Handler struct {
	svc    Service
	cfg    *config.Provider
	logger *zap.Logger
	fbAuth *middleware.FirebaseAuth
}

// NewHandler creates a new user handler.
func NewHandler(svc Service, cfg *config.Provider, logger *zap.Logger, fbAuth *middleware.FirebaseAuth) *Handler {
	return &Handler{svc: svc, cfg: cfg, logger: logger, fbAuth: fbAuth}
}

// RegisterRoutes mounts user routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users")

	users.Post("/sync", h.Upsert)
	users.Get("", h.List)
	users.Get("/:id", h.GetByID)

	var authMiddleware fiber.Handler
	if h.cfg.Auth.Mode == "jwt" {
		authMiddleware = middleware.JWTAuth(h.cfg.Auth.JWTSecret)
	} else {
		if h.fbAuth != nil {
			authMiddleware = h.fbAuth.Auth()
		} else {
			authMiddleware = func(c *fiber.Ctx) error { return c.Next() }
		}
	}

	users.Patch("/:id", authMiddleware, h.UpdateProfile)
	users.Delete("/:id", authMiddleware, h.Delete)
}

// Upsert godoc
// @Summary Sync/upsert a user from Firebase
// @Tags users
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "User data"
// @Success 200 {object} utils.Envelope{data=UserResponse}
// @Failure 400 {object} utils.Envelope
// @Router /users/sync [post]
func (h *Handler) Upsert(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "invalid request body")
	}

	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return utils.ResponseValidationError(c, errs)
	}

	resp, err := h.svc.UpsertByFirebaseUID(c.Context(), req)
	if err != nil {
		return utils.ParseErrHTTP(c, err, nil, h.logger)
	}

	return utils.ResponseOK(c, resp)
}

// GetByID godoc
// @Summary Get a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.Envelope{data=UserResponse}
// @Failure 404 {object} utils.Envelope
// @Router /users/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return utils.ParseErrHTTP(c, err, nil, h.logger)
	}

	return utils.ResponseOK(c, resp)
}

// List godoc
// @Summary List users with pagination
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Param sort query string false "Sort field (prefix with - for DESC)"
// @Success 200 {object} utils.Envelope{data=UserListResponse}
// @Router /users [get]
func (h *Handler) List(c *fiber.Ctx) error {
	page, pageSize := utils.ParsePaginationParams(c)

	req := ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   c.Query("search"),
		Sort:     c.Query("sort"),
	}

	resp, err := h.svc.List(c.Context(), req)
	if err != nil {
		return utils.ParseErrHTTP(c, err, nil, h.logger)
	}

	return utils.ResponseOKWithMeta(c, resp.Users, resp.Pagination)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body UpdateUserRequest true "Profile data"
// @Success 200 {object} utils.Envelope{data=UserResponse}
// @Failure 400 {object} utils.Envelope
// @Failure 401 {object} utils.Envelope
// @Security BearerAuth
// @Router /users/{id} [patch]
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "invalid request body")
	}

	if errs := utils.ValidateStruct(req); len(errs) > 0 {
		return utils.ResponseValidationError(c, errs)
	}

	resp, err := h.svc.UpdateProfile(c.Context(), id, req)
	if err != nil {
		return utils.ParseErrHTTP(c, err, nil, h.logger)
	}

	return utils.ResponseOK(c, resp)
}

// Delete godoc
// @Summary Delete a user
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.Envelope
// @Failure 401 {object} utils.Envelope
// @Failure 404 {object} utils.Envelope
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return utils.ParseErrHTTP(c, err, nil, h.logger)
	}

	return utils.ResponseOK(c, fiber.Map{"message": "user deleted"})
}
