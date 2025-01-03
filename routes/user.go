package routes

import (
	"myIris/controllers"
	"myIris/middlewares"

	"github.com/kataras/iris/v12"
)

func UserRoutes(app *iris.Application) {
	user := app.Party("/user")
	{
		user.Use(middlewares.AuthMiddleware())
		user.Get("/logout", controllers.Logout)
		// user.Get("/profile", controllers.GetProfile)
		// user.Put("/profile", controllers.UpdateProfile)
		// user.Get("/profile/{userId}", controllers.GetProfile)
		// user.Get("/profile/{userId}/posts", controllers.GetPostsByUser)
		// user.Get("/profile/{userId}/posts/{postId}", controllers.GetPostByUser)
		// user.Post("/profile/{userId}/posts", controllers.CreatePost)
		// user.Put("/profile/{userId}/posts/{postId}", controllers.UpdatePost)
		// user.Delete("/profile/{userId}/posts/{postId}", controllers.DeletePost)
	}

}
