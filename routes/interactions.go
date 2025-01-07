package routes

import (
	"myIris/controllers"
	"myIris/middlewares"

	"github.com/kataras/iris/v12"
)

func InteractionsRoutes(app *iris.Application) {
	interactions := app.Party("/interactions")
	{
		interactions.Use(middlewares.AuthMiddleware())
		interactions.Post("/{lid}", controllers.CreateInteractionWithLead)
		interactions.Put("/{iid}", controllers.UpdateNoteOfInteraction)
		interactions.Get("/", controllers.GetInteractionHistory)
	}
}
