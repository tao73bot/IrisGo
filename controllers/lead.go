package controllers

import (
	"myIris/db"
	"myIris/models"
	"myIris/utils"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func CreateLead(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}

	tokenString := ctx.GetHeader("Authorization")
	claims, err := utils.ValidateTokenIris(tokenString)
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	var body struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Sources string `json:"sources"`
		Status  string `json:"status"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(err.Error()))
		return
	}
	lead := models.Lead{
		UserID: userID,
		Name:   body.Name,
		Email:  body.Email,
		Phone:  body.Phone,
		Source: body.Sources,
		Status: body.Status,
	}
	result := db.DB.Create(&lead)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"message": "Lead created successfully"})
}

func GetAllLeads(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "admin"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	tokenString := ctx.GetHeader("Authorization")
	_, err := utils.ValidateTokenIris(tokenString)
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	var leads []models.Lead
	result := db.DB.Find(&leads)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"leads": leads})
}

func GetAllLeadByUser(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	tokenString := ctx.GetHeader("Authorization")
	claims, err := utils.ValidateTokenIris(tokenString)
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	var leads []models.Lead
	result := db.DB.Where("user_id = ?", claims.UserID).Find(&leads)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"leads": leads})
}

func GetLeadByName(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	tokenString := ctx.GetHeader("Authorization")
	_, err := utils.ValidateTokenIris(tokenString)
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	name := ctx.Params().Get("name")
	var leads []models.Lead
	result := db.DB.Where("name = ?", name).Find(&leads)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"leads": leads})
}

func GetLeadByID(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	tokenString := ctx.GetHeader("Authorization")
	_, err := utils.ValidateTokenIris(tokenString)
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	id, err := uuid.Parse(ctx.Params().Get("id"))
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	var lead models.Lead
	result := db.DB.Where("lead_id = ?", id).First(&lead)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"lead": lead})
}

func UpdateLeadInfo(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	tokenString := ctx.GetHeader("Authorization")
	claims, err := utils.ValidateTokenIris(tokenString)
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	var body struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Sources string `json:"sources"`
		Status  string `json:"status"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	id, err := uuid.Parse(ctx.Params().Get("id"))
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	var lead models.Lead
	result := db.DB.Where("lead_id = ?", id).First(&lead)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}

	if lead.UserID.String() != claims.UserID {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("You are not authorized to update this lead"))
		return
	}
	if body.Name != "" {
		lead.Name = body.Name
	}
	if body.Email != "" {
		lead.Email = body.Email
	}
	if body.Phone != "" {
		lead.Phone = body.Phone
	}
	if body.Sources != "" {
		lead.Source = body.Sources
	}
	if body.Status != "" {
		lead.Status = body.Status
	}
	result = db.DB.Save(&lead)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"message": "Lead updated successfully"})
}

func DeleteLead(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	tokenString := ctx.GetHeader("Authorization")
	claims, err := utils.ValidateTokenIris(tokenString)
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	id, err := uuid.Parse(ctx.Params().Get("id"))
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	var lead models.Lead
	result := db.DB.Where("lead_id = ?", id).First(&lead)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	if lead.UserID.String() != claims.UserID {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("You are not authorized to delete this lead"))
		return
	}
	result = db.DB.Delete(&lead)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"message": "Lead deleted successfully"})
}
