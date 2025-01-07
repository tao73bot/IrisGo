package controllers

import (
	"myIris/db"
	"myIris/models"
	"myIris/utils"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func CreateCustomer(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
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
	uid := uuid.MustParse(claims.UserID)
	var body struct {
		Address     string `json:"address"`
		CompanyName string `json:"company_name"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	var lead models.Lead
	result := db.DB.Where("lead_id = ? and user_id = ?", lid, uid).First(&lead)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
			Title("Not Found").
			Detail("Lead not found or you are not authorized to create a customer for this lead"))
		return
	}
	if lead.Status != "qualified" {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
			Detail("Lead not qualified"))
		return
	}
	customer := models.Customer{
		LeadID:      lid,
		UserID:      uid,
		Address:     body.Address,
		CompanyName: body.CompanyName,
	}
	lidCnt := db.DB.Where("lead_id = ?", lid).Find(&models.Customer{}).RowsAffected
	if lidCnt > 0 {
		ctx.StopWithProblem(iris.StatusConflict, iris.NewProblem().
			Title("Conflict").
			Detail("Customer already exists"))
		return
	}
	result = db.DB.Create(&customer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	// join lead and customer
	var updatedCustomer struct {
		CustomerID  uuid.UUID `json:"customer_id"`
		LeadID      uuid.UUID `json:"lead_id"`
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		Phone       string    `json:"phone"`
		Status      string    `json:"status"`
		Source      string    `json:"source"`
		Address     string    `json:"address"`
		CompanyName string    `json:"company_name"`
	}
	x := db.DB.Table("customers").
		Select("customers.*, leads.*").
		Joins("JOIN leads ON leads.lead_id = customers.lead_id").
		Where("customers.customer_id = ?", customer.CustomerID).
		Scan(&updatedCustomer)
	if x.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(x.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{
		"message":  "Customer created successfully",
		"customer": updatedCustomer,
	})
}

func GetAllCustomers(ctx iris.Context) {
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
	var updatedCustomer []struct {
		CustomerID  uuid.UUID `json:"customer_id"`
		LeadID      uuid.UUID `json:"lead_id"`
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		Phone       string    `json:"phone"`
		Status      string    `json:"status"`
		Source      string    `json:"source"`
		Address     string    `json:"address"`
		CompanyName string    `json:"company_name"`
	}
	result := db.DB.Table("customers").
		Select("customers.*, leads.*").
		Joins("JOIN leads ON leads.lead_id = customers.lead_id").
		Scan(&updatedCustomer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"customers": updatedCustomer})
}

func GetCustomerByID(ctx iris.Context) {
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
	cid := uuid.MustParse(ctx.Params().Get("cid"))
	uid := uuid.MustParse(claims.UserID)
	var updatedCustomer struct {
		CustomerID  uuid.UUID `json:"customer_id"`
		LeadID      uuid.UUID `json:"lead_id"`
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		Phone       string    `json:"phone"`
		Status      string    `json:"status"`
		Source      string    `json:"source"`
		Address     string    `json:"address"`
		CompanyName string    `json:"company_name"`
	}
	result := db.DB.Table("customers").
		Select("customers.*, leads.*").
		Joins("JOIN leads ON leads.lead_id = customers.lead_id").
		Where("customers.customer_id = ? AND customers.user_id = ?", cid, uid).
		Scan(&updatedCustomer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"customer": updatedCustomer})
}

func GetCustomersOfUser(ctx iris.Context) {
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
	uid := uuid.MustParse(claims.UserID)
	var updatedCustomer []struct {
		CustomerID  uuid.UUID `json:"customer_id"`
		LeadID      uuid.UUID `json:"lead_id"`
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		Phone       string    `json:"phone"`
		Status      string    `json:"status"`
		Source      string    `json:"source"`
		Address     string    `json:"address"`
		CompanyName string    `json:"company_name"`
	}
	result := db.DB.Table("customers").
		Select("customers.*, leads.*").
		Joins("JOIN leads ON leads.lead_id = customers.lead_id").
		Where("customers.user_id = ?", uid).
		Scan(&updatedCustomer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"customers": updatedCustomer})
}

func GetCustomersByUserID(ctx iris.Context) {
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
	uid := uuid.MustParse(ctx.Params().Get("uid"))
	var updatedCustomer []struct {
		CustomerID  uuid.UUID `json:"customer_id"`
		LeadID      uuid.UUID `json:"lead_id"`
		UserID      uuid.UUID `json:"user_id"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		Phone       string    `json:"phone"`
		Status      string    `json:"status"`
		Source      string    `json:"source"`
		Address     string    `json:"address"`
		CompanyName string    `json:"company_name"`
	}
	result := db.DB.Table("customers").
		Select("customers.*, leads.*").
		Joins("JOIN leads ON leads.lead_id = customers.lead_id").
		Where("customers.user_id = ?", uid).
		Scan(&updatedCustomer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"customers": updatedCustomer})
}

func UpdateCustomerInfo(ctx iris.Context) {
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
	cid := uuid.MustParse(ctx.Params().Get("cid"))
	uid := uuid.MustParse(claims.UserID)
	var body struct {
		Address     string `json:"address"`
		CompanyName string `json:"company_name"`
	}
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Bad Request").
			Detail(err.Error()))
		return
	}
	var customer models.Customer
	result := db.DB.Where("customer_id = ? AND user_id = ?", cid, uid).First(&customer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
			Title("Not Found").
			Detail("Customer not found or you are not authorized to update this customer"))
		return
	}
	if body.Address != "" {
		customer.Address = body.Address
	}
	if body.CompanyName != "" {
		customer.CompanyName = body.CompanyName
	}
	result = db.DB.Save(&customer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"message": "Customer updated successfully", "customer": customer})
}

func DeleteCustomer(ctx iris.Context) {
	if err := utils.CheckUserRoles(ctx, "user"); err != nil {
		ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
			Title("Forbidden").
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
	cid := uuid.MustParse(ctx.Params().Get("cid"))
	uid := uuid.MustParse(claims.UserID)
	var customer models.Customer
	result := db.DB.Where("customer_id = ? AND user_id = ?", cid, uid).First(&customer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
			Title("Not Found").
			Detail("Customer not found or you are not authorized to delete this customer"))
		return
	}
	result = db.DB.Delete(&customer)
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error").
			Detail(result.Error.Error()))
		return
	}
	ctx.JSON(iris.Map{"message": "Customer deleted successfully"})
}
