package handlers

import (
	"go-rest-template/internal/models"
	"go-rest-template/internal/services"
	"go-rest-template/pkg/API"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: *service,
	}
}

// GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		API.Error(c, "failed to list users")
		return
	}
	API.Success(c, "success", users)
}

// GET /users/me
func (h *UserHandler) GetMyDetails(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		API.Unauthorized(c, "unauthorized")
		return
	}
	currentUser := user.(*models.User)

	u, err := h.service.GetUserByEmail(c.Request.Context(), currentUser.Email)
	if err != nil {
		API.Error(c, "failed to get user details")
		return
	}
	API.Success(c, "success", u)
}
