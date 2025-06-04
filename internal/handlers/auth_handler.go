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

	if exists, _ := h.service.GetUserByEmail(c.Request.Context(), u.Email); exists != nil {
		API.Error(c, "user already exists")
		return
	}
	if err := h.service.CreateUser(c.Request.Context(), &u); err != nil {
		API.InternalServerError(c, "failed to create user")
		return
	}
	u.Password = ""
	API.Success(c, "user created successfully", u)
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var u dto.User_LoginRequest

	if err := c.ShouldBindJSON(&u); err != nil {
		errors := pkg.TagValidationErrors(err, &u)
		API.ValidationsErrors(c, errors)
		return
	}

	exists, err := h.service.Authenticate(c.Request.Context(), u.Email, u.Password)
	if err != nil {
		API.Error(c, err.Error())
		return
	}

	token, err := pkg.GenerateToken(u.Email)
	if err != nil {
		API.InternalServerError(c, "failed to generate token")

	}
	API.SendJWTtoken(c, token, "logged in successfully", map[string]any{"token": token, "user": exists})
}
