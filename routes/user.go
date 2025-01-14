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
		user.Get("/", controllers.GetUsers)
		user.Get("/get_access_token", controllers.GenerateAccessTokenUsingRefreshToken)
		user.Get("/{userId:string}", controllers.GetUser)
		user.Get("/{userId:string}/another", controllers.GetAnotherUser)
		user.Put("/{userId:string}", controllers.UpdateUser)
		user.Put("/{userId:string}/password", controllers.UpdateUserPassword)
		user.Put("/{userId:string}/role", controllers.UpdateUserRole)
		user.Delete("/{userId:string}", controllers.DeleteUser)
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
