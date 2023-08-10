package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/alvinahb/supply-me/db/sqlc"
	"github.com/gin-gonic/gin"
)

type orderRequest struct {
	FromCompanyID int64 `json:"from_company_id" binding:"required,min=1"`
	ToCompanyID   int64 `json:"to_company_id" binding:"required,min=1"`
	ProductID     int64 `json:"product_id" binding:"required,min=1"`
	Amount        int32 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req orderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validProvider(ctx, req.FromCompanyID, req.ProductID, req.Amount) {
		return
	}

	arg := db.OrderTxParams{
		FromCompanyID: req.FromCompanyID,
		ToCompanyID:   req.ToCompanyID,
		ProductID:     req.ProductID,
		Amount:        req.Amount,
	}

	result, err := server.store.OrderTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validProvider(ctx *gin.Context, providerID int64, productID int64, amount int32) bool {
	provider, err := server.store.GetCompany(ctx, providerID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	// Check if order targets a provider
	if provider.CompanyType != "Fournisseur" {
		err = fmt.Errorf("Company %s is not a Provider but a %s",
			provider.CompanyName, provider.CompanyType)
		ctx.JSON(http.StatusBadRequest, err)
		return false

		// Check if provider has enough amount of order's product
	} else {
		arg := db.GetCompanyProductInventoryParams{
			CompanyID: providerID,
			ProductID: productID,
		}

		inventory, err := server.store.GetCompanyProductInventory(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return false
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return false
		}

		if inventory.AmountAvailable < amount {
			err = fmt.Errorf("Ordered %d but company %s only has %d of this product",
				amount, provider.CompanyName, inventory.AmountAvailable)
			ctx.JSON(http.StatusBadRequest, err)
			return false
		}
	}

	return true
}
