package controllers

import (
	"log"
	"myIris/db"
	"myIris/models"
	"myIris/utils"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
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
	signer := jwt.NewSigner(jwt.HS256, []byte(os.Getenv("JWT_SECRET")), 20*time.Minute)
	refreshSigner := jwt.NewSigner(jwt.HS256, []byte(os.Getenv("JWT_SECRET")), 1*time.Hour)
	token, err := utils.GenerateTokenIris(signer, user.Email, user.Name, user.Role, user.UserID.String())
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	refreshToken, err := utils.GenerateTokenIris(refreshSigner, user.Email, user.Name, user.Role, user.UserID.String())
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	// token, refreshToken, err := utils.GenerateAllTokens(user.Email, user.Name, user.Role, user.UserID.String())
	// if err != nil {
	// 	ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
	// 		Title("Internal Server Error").
	// 		Detail(err.Error()))
	// 	return
	// }
	// ctx.SetCookie(&http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refreshToken,
	// 	HttpOnly: true,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(1 * time.Hour),
	// })
	ctx.SetCookieKV("refresh_token", refreshToken, iris.CookiePath("/"), iris.CookieExpires(1*time.Hour), iris.CookieHTTPOnly(true))
	ctx.JSON(iris.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func Logout(ctx iris.Context) {
	accessToken := ctx.GetHeader("Authorization")
	tb, err := utils.NewTokenBlocklist(os.Getenv("REDIS_HOST"))
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	err = tb.InvalidateTokenIris(accessToken, 20*time.Minute)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	log.Println("pass Access")
	refreshToken := ctx.GetCookie("refresh_token")
	if refreshToken == "" {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Request").
			Detail("Refresh token not found"))
		return
	}
	log.Println(refreshToken)
	err = tb.InvalidateTokenIris(refreshToken, 1*time.Hour)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
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
		"users":   users,
	})
}

func GetUser(ctx iris.Context) {
	uid := uuid.MustParse(ctx.Params().Get("userId"))
	if err := utils.MatchRoleToUid(ctx, uid); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail(err.Error()))
		return
	}
	var user models.User
	result := db.DB.Where("user_id = ?", uid).First(&user)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Request").
			Detail("User not found"))
		return
	}
	ctx.JSON(iris.Map{
		"message": "User retrieved successfully",
		"user":    user,
	})
}

func GetAnotherUser(ctx iris.Context) {
	uid := uuid.MustParse(ctx.Params().Get("userId"))
	var user models.User
	result := db.DB.Where("user_id = ?", uid).First(&user)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Request").
			Detail("User not found"))
		return
	}
	ctx.JSON(iris.Map{
		"message": "User retrieved successfully",
		"name":   user.Name,
		"email":    user.Email,
		"role":    user.Role,
	})
}

func UpdateUser(ctx iris.Context) {
	uid := uuid.MustParse(ctx.Params().Get("userId"))
	if err := utils.MatchRoleToUid(ctx, uid); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail(err.Error()))
		return
	}
	var user models.User
	result := db.DB.Where("user_id = ?", uid).First(&user)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Request").
			Detail("User not found"))
		return
	}
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail(err.Error()))
		return
	}
	user.Name = body.Name
	user.Email = body.Email
	result = db.DB.Save(&user)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message": "User updated successfully",
		"user":    user,
	})
}

func UpdateUserPassword(ctx iris.Context) {
	uid := uuid.MustParse(ctx.Params().Get("userId"))
	if err := utils.MatchRoleToUid(ctx, uid); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail(err.Error()))
		return
	}
	var user models.User
	result := db.DB.Where("user_id = ?", uid).First(&user)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Request").
			Detail("User not found"))
		return
	}
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail(err.Error()))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.CurrentPassword))
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail("Invalid Password"))
		return
	}
	if body.NewPassword != body.ConfirmPassword {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail("Passwords do not match"))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	user.Password = string(hash)
	result = db.DB.Save(&user)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message": "Password updated successfully",
	})
}

func UpdateUserRole(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "admin"); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail(err.Error()))
		return
	}
	uid := uuid.MustParse(ctx.Params().Get("userId"))
	var user models.User
	result := db.DB.Where("user_id = ?", uid).First(&user)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Request").
			Detail("User not found"))
		return
	}
	var body struct {
		Role string `json:"role"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Credentials").
			Detail(err.Error()))
		return
	}
	user.Role = body.Role
	result = db.DB.Save(&user)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message": "Role updated successfully",
		"user":    user,
	})
}

func DeleteUser(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "admin"); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail(err.Error()))
		return
	}
	uid := uuid.MustParse(ctx.Params().Get("userId"))
	var user models.User
	result := db.DB.Where("user_id = ?", uid).First(&user)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid Request").
			Detail("User not found"))
		return
	}
	result = db.DB.Delete(&user)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message": "User deleted successfully",
	})
}