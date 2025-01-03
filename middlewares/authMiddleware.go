package middlewares

import (
	"myIris/utils"

	"github.com/kataras/iris/v12"
)

func AuthMiddleware() iris.Handler {
	return func(ctx iris.Context) {
		clientToken := ctx.GetHeader("Authorization")
		if clientToken == "" {
			ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
				Title("Unauthorized").
				Detail("No token provided"))
			return
		}
		claims, err := utils.ValidateToken(clientToken)
		if err != nil {
			ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
				Title("Unauthorized").
				Detail(err.Error()))
			return
		}
		for _, v := range utils.BlockList {
			if v == clientToken {
				ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
					Title("Unauthorized").
					Detail("Token is invalid"))
				return
			}
		}
		ctx.Values().Set("email", claims.Email)
		ctx.Values().Set("role", claims.Role)
		ctx.Values().Set("userID", claims.UserID)
		ctx.Next()
	}
}
