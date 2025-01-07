package main

import (
	"fmt"
	"myIris/db"
	"myIris/routes"
	"os"

	"github.com/kataras/iris/v12"
)

func init() {
	db.LoadEnv()
	db.ConnectDB()
	db.Migrate()
}

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello World!"})
	})
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.LeadRoutes(app)
	routes.CustomerRoutes(app)
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	fmt.Println("Server running on port: " + Port)
	app.Listen(":" + Port)
}
