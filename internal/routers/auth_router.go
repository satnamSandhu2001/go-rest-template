package routers

import (
	"go-rest-template/internal/handlers"
	"go-rest-template/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterAuthRouter(rg *gin.RouterGroup, db *sqlx.DB) {
	s := services.NewUserService(db)
	h := handlers.NewAuthHandler(s)

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/signup", h.Signup)
		authGroup.POST("/login", h.Login)
	}
}
