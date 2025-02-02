package routes

import (
	"myIris/controllers"

	"github.com/kataras/iris/v12"
)

func AuthRoutes(app *iris.Application) {
	auth := app.Party("/auth")
	{
		auth.Post("/signup", controllers.SignUp)
		auth.Post("/signin", controllers.Login)
		auth.Get("/verify/{token:string}", controllers.VerifyEmail)
		auth.Post("/forgot-password", controllers.ForgotPassword)
		auth.Put("/reset-password/{token:string}", controllers.ResetPassword)
	}
}
