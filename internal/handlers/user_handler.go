package handlers

import (
	"go-rest-template/internal/db"
	"go-rest-template/internal/dto"
	"go-rest-template/internal/middlewares"
	"go-rest-template/internal/repository"
	"go-rest-template/pkg"
	"go-rest-template/pkg/api"
	"go-rest-template/pkg/logger"
	"net/http"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// POST /auth/signup
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	body, validationErrors, err := pkg.BindAndValidate[dto.User_Signup_Request](r)
	if err != nil {
		api.Error(w, err.Error())
		return
	}
	if validationErrors != nil {
		api.ValidationErrors(w, validationErrors)
		return
	}

	existingUser, err := h.repo.FindByEmail(r.Context(), body.Email)
	if err != nil {
		logger.Debug(err)
		api.InternalServerError(w)
		return
	}
	if existingUser.ID != 0 {
		api.Error(w, "User already exists")
		return
	}
	hash, err := pkg.GenerateHash(body.Password)
	if err != nil {
		logger.Error(err)
		api.InternalServerError(w)
		return
	}
	id, err := h.repo.CreateNew(r.Context(), db.CreateUserParams{
		Email:        body.Email,
		PasswordHash: hash,
	})
	if err != nil {
		logger.Error(err)
		api.InternalServerError(w)
		return
	}
	newUser, err := h.repo.FindById(r.Context(), id)
	if err != nil {
		logger.Error(err)
		api.InternalServerError(w)
		return
	}
	if newUser.ID == 0 {
		api.Error(w, "User not found")
		return
	}

	token, err := pkg.GenerateToken(newUser.Email)
	if err != nil {
		logger.Error(err)
		api.InternalServerError(w)
		return
	}
	api.SendJWTtoken(w, r, token, "Signed up successfully", newUser)
}

// POST /auth/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, validationErrors, err := pkg.BindAndValidate[dto.User_Login_Request](r)
	if err != nil {
		api.Error(w, err.Error())
		return
	}
	if validationErrors != nil {
		api.ValidationErrors(w, validationErrors)
		return
	}

	user, err := h.repo.FindByEmail(r.Context(), body.Email)
	if err != nil {
		logger.Error(err)
		api.InternalServerError(w)
		return
	}
	if user.ID == 0 {
		api.Error(w, "Invalid credentials")
		return
	}
	if notMatched := pkg.CompareHashAndPassword(user.PasswordHash, body.Password); notMatched != nil {
		api.Error(w, "Invalid credentials")
		return
	}

	token, err := pkg.GenerateToken(user.Email)
	if err != nil {
		logger.Error(err)
		api.InternalServerError(w)
		return
	}
	api.SendJWTtoken(w, r, token, "Logged in successfully", user)
}

// GET /auth/logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	api.Logout(w)
}

// GET /user/profile
func (h *UserHandler) GetMyDetails(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	api.Success(w, "Success", user)
}
