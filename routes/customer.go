package routes

import (
	"myIris/controllers"
	"myIris/middlewares"

	"github.com/kataras/iris/v12"
)

func CustomerRoutes(app *iris.Application) {
	customer := app.Party("/customer")
	{
		customer.Use(middlewares.AuthMiddleware())
		customer.Post("/{lid}", controllers.CreateCustomer)
		customer.Get("/", controllers.GetAllCustomers)
		customer.Get("/{cid}", controllers.GetCustomerByID)
		customer.Get("/user",controllers.GetCustomersOfUser)
		customer.Get("/user/{uid}",controllers.GetCustomersByUserID)
		customer.Put("/{cid}", controllers.UpdateCustomerInfo)
		customer.Delete("/{cid}", controllers.DeleteCustomer)
	}
}
