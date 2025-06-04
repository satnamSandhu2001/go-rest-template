package routers

import (
	"go-rest-template/internal/handlers"
	"go-rest-template/internal/middlewares"
	"go-rest-template/internal/services"

	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"
)

func RegisterUserRouter(rg *gin.RouterGroup, db *sqlx.DB) {

	s := services.NewUserService(db)
	h := handlers.NewUserHandler(s)

	usersGroup := rg.Group("/users")
	authGroup := usersGroup.Group("", middlewares.AuthMiddleware(s))
	{
		authGroup.GET("me", h.GetMyDetails)
		authGroup.GET("", h.ListUsers)
	}

}
