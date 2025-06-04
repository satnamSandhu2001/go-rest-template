package handlers

import (
	"go-rest-template/internal/dto"
	"go-rest-template/internal/services"
	"go-rest-template/pkg"
	"go-rest-template/pkg/API"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{
		service: *service,
	}
}

// POST /auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var u dto.User_RegisterRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		errors := pkg.TagValidationErrors(err, &u)
		API.ValidationsErrors(c, errors)
		return
	}
	userExists, err := h.service.GetUserByEmail(c.Request.Context(), u.Email)
	if err != nil {
		API.InternalServerError(c, "Failed to signup", err)
		return
	}
	if userExists != nil {
		API.Error(c, "User with this email already exists")
		return
	}
	if err := h.service.CreateUser(c.Request.Context(), &u); err != nil {
		API.InternalServerError(c, "Failed to signup", err)
		return
	}

	newUser, err := h.service.GetUserByID(c.Request.Context(), u.ID)
	if (err != nil) || (newUser == nil) {
		API.InternalServerError(c, "Failed to signup", err)
		return
	}

	token, err := pkg.GenerateToken(u.Email)
	if err != nil {
		API.InternalServerError(c, "Failed to generate token", err)

	}
	API.SendJWTtoken(c, token, "Account Created Successfully", map[string]any{"token": token, "user": newUser})
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var u dto.User_LoginRequest

	if err := c.ShouldBindJSON(&u); err != nil {
		errors := pkg.TagValidationErrors(err, &u)
		API.ValidationsErrors(c, errors)
		return
	}

	user, err := h.service.Authenticate(c.Request.Context(), u.Email, u.Password)
	if err != nil {
		API.InternalServerError(c, "invalid credentials", err)
		return
	}
	if user == nil {
		API.Error(c, "invalid credentials")
		return
	}

	token, err := pkg.GenerateToken(u.Email)
	if err != nil {
		API.InternalServerError(c, "failed to generate token", err)
	}

	API.SendJWTtoken(c, token, "logged in successfully", map[string]any{"token": token, "user": user})
}
