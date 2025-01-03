package utils

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func CheckUserRoles(ctx iris.Context, role string) (err error) {
	userRole := ctx.Values().GetString("role")
	if userRole != role {
		err = errors.New("you are not authorized to access this route")
		return err
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
		err = errors.New("you are not authorized to access this route")
		return err
	}
	err = CheckUserRoles(ctx, userRole)
	return err
}
