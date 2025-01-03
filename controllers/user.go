package controllers

import (
	"myIris/db"
	"myIris/models"

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
		"user": user,
	})
}
