package controllers

import (
	"myIris/db"
	"myIris/models"
	"myIris/utils"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func CreateInteractionWithLead(ctx iris.Context) {
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
	lid := uuid.MustParse(ctx.Params().Get("lid"))
	var lead models.Lead
	result := db.DB.Where("lead_id = ?", lid).First(&lead)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
			Title("Not Found").
			Detail("Lead not found"))
		return
	}
	uid := uuid.MustParse(claims.UserID)
	var body struct {
		Type  string `json:"type"`
		Notes string `json:"notes"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	interaction := models.Interaction{
		LeadID: lid,
		UserID: uid,
		Type:   body.Type,
		Notes:  body.Notes,
	}
	result = db.DB.Create(&interaction)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message":     "Interaction created successfully",
		"interaction": interaction,
	})
}

func UpdateNoteOfInteraction(ctx iris.Context) {
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
	iid := uuid.MustParse(ctx.Params().Get("iid"))
	var interaction models.Interaction
	result := db.DB.Where("interaction_id = ?", iid).First(&interaction)
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
			Title("Not Found").
			Detail("Interaction not found"))
		return
	}
	uid := uuid.MustParse(claims.UserID)
	if interaction.UserID != uid {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("You are not authorized to update this interaction"))
		return
	}
	var body struct {
		Notes string `json:"notes"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	interaction.Notes = body.Notes
	result = db.DB.Save(&interaction)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message":     "Interaction updated successfully",
		"interaction": interaction,
	})
}

func GetInteractionHistory(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "admin"); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	_, err := utils.ValidateTokenIris(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail(err.Error()))
		return
	}
	interactions := []models.Interaction{}
	result := db.DB.Find(&interactions)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"interactions": interactions,
	})
}
