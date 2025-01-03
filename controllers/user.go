package controllers

import (
	"myIris/db"
	"myIris/models"
	"myIris/utils"
	"net/http"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx iris.Context) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail(err.Error()))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
	}
	x := db.DB.Where("email = ?", body.Email).First(&user).RowsAffected
	if x > 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail("Email already exists"))
		return
	}
	result := db.DB.Create(&user)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

func Login(ctx iris.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail(err.Error()))
		return
	}
	var user models.User
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail("User not found"))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail("Invalid Email or Password"))
		return
	}
	token, refreshToken, err := utils.GenerateAllTokens(user.Email, user.Name, user.Role, user.UserID.String())
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	ctx.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
	})
	ctx.JSON(iris.Map{
		"message":      "Login successful",
		"token":        token,
	})
}


func Logout(ctx iris.Context) {
	accessToken := ctx.GetHeader("Authorization")
	utils.InvalidateToken(accessToken)
	refreshToken := ctx.GetCookie("refresh_token")
	utils.InvalidateToken(refreshToken)
	ctx.RemoveCookie("refresh_token")
	ctx.JSON(iris.Map{
		"message": "Logout successful",
	})
}

func GetUsers(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "admin"); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail(err.Error()))
		return
	}
	var users []models.User
	result := db.DB.Find(&users)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message": "Users retrieved successfully",
		"users": users,
	})
}