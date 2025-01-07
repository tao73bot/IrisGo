package routes

import (
	"myIris/controllers"
	"myIris/middlewares"

	"github.com/kataras/iris/v12"
)

func LeadRoutes(app *iris.Application) {
	leadRoutes := app.Party("/leads")
	{
		leadRoutes.Use(middlewares.AuthMiddleware())
		leadRoutes.Post("/", controllers.CreateLead)
		leadRoutes.Get("/", controllers.GetAllLeads)
		leadRoutes.Get("/{id}", controllers.GetLeadByID)
		leadRoutes.Get("/user", controllers.GetAllLeadByUser)
		leadRoutes.Get("/get_lead_by_name/{name}", controllers.GetLeadByName)
		leadRoutes.Patch("/{id}", controllers.UpdateLeadInfo)
		leadRoutes.Delete("/{id}", controllers.DeleteLead)
	}
}
