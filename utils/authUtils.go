package utils

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func CheckUserRoles(ctx iris.Context, role string) (err error) {
	userRole := ctx.Values().GetString("role")
	if userRole != role {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail("You are not authorized to access this route"))
		return
	}
	return nil
}

func MatchRoleToUid(ctx iris.Context, userId uuid.UUID) (err error) {
	userRole := ctx.Values().GetString("role")
	uid, err := uuid.Parse(ctx.Params().Get("userId"))
	if err != nil {
		return err
	}
	if userRole == "user" && userId != uid {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail("You are not authorized to access this route"))
		return
	}
	err = CheckUserRoles(ctx, userRole)
	return err
}
