package main

import (
	"fmt"
	"go-rest-template/database"
	"go-rest-template/internal/routers"
	"go-rest-template/pkg"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	conn := database.Connect()
	database.RunMigrations(conn)

	r := routers.InitRouter(conn)
	r.SetTrustedProxies(nil)

	if err := r.Run(fmt.Sprintf(":%v", pkg.Config().PORT)); err != nil {
		panic(err)
	}
}
