package main

import (
	"fmt"
	"myIris/db"
	"myIris/routes"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
)

func init() {
	db.LoadEnv()
	db.ConnectDB()
	db.Migrate()
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]

func main() {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello World!"})
	})
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.LeadRoutes(app)
	routes.CustomerRoutes(app)
	routes.InteractionsRoutes(app)

	app.HandleDir("/docs", "./docs", iris.DirOptions{IndexName: "index.html"})
	app.Get("/swagger/{any:path}", swagger.WrapHandler(swaggerFiles.Handler, swagger.URL("/docs/swagger.json")))
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	fmt.Println("Server running on port: " + Port)
	app.Listen(":" + Port)
}
