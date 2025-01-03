package routes

import (
	"myIris/controllers"

	"github.com/kataras/iris/v12"
)

func AuthRoutes(app *iris.Application) {
	auth := app.Party("/auth")
	{
		auth.Post("/signup", controllers.SignUp)
		// auth.Post("/login", login)
	}
}