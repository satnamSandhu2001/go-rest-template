package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitRouter(db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	RegisterAuthRouter(v1, db)
	RegisterUserRouter(v1, db)

	return r

}
