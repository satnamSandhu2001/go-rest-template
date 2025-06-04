package middlewares

import (
	"go-rest-template/internal/services"
	"go-rest-template/pkg"
	"go-rest-template/pkg/API"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("Authorization")
		if err != nil {
			API.Unauthorized(ctx, "unauthorized")
			return
		}

		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			API.Unauthorized(ctx, "unauthorized")
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		email, err := pkg.ValidateToken(token)
		if err != nil {
			API.Unauthorized(ctx, "unauthorized")
			return
		}

		user, err := userService.GetUserByEmail(ctx.Request.Context(), email)
		if (user == nil) || (err != nil) {
			API.Unauthorized(ctx, "unauthorized")
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
