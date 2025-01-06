package middlewares

import (
	"log"
	"myIris/utils"
	"os"

	"github.com/kataras/iris/v12"
)

func AuthMiddleware() iris.Handler {
	return func(ctx iris.Context) {
		clientToken := ctx.GetHeader("Authorization")
		if clientToken == "" {
			ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
				Title("Unauthorized").
				Detail("No token provided"))
			ctx.StopExecution()
			return
		}
		claims, err := utils.ValidateTokenIris(clientToken)
		if err != nil {
			ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
				Title("Unauthorized").
				Detail(err.Error()))
			ctx.StopExecution()
			return
		}
		tb, err := utils.NewTokenBlocklist(os.Getenv("REDIS_HOST"))
		if err != nil {
			log.Println("sihsi")
			ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
				Title("Internal Server Error").
				Detail(err.Error()))
			ctx.StopExecution()
			return
		}
		if tb.IsTokenInvalidIris(clientToken) {
			ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
				Title("Unauthorized").
				Detail("Token is blocked"))
			ctx.StopExecution()
			return
		}

		ctx.Values().Set("email", claims.Email)
		ctx.Values().Set("role", claims.Role)
		ctx.Values().Set("userID", claims.UserID)
		ctx.Next()
	}
}
