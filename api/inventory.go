package api

import (
	"database/sql"
	"net/http"

	db "github.com/alvinahb/supply-me/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createInventoryRequest struct {
	CompanyID       int64 `json:"company_id" binding:"required,min=1"`
	ProductID       int64 `json:"product_id" binding:"required,min=1"`
	AmountAvailable int32 `json:"amount_available" binding:"required,min=0"`
}

func (server *Server) createInventory(ctx *gin.Context) {
	var req createInventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateInventoryParams{
		CompanyID:       req.CompanyID,
		ProductID:       req.ProductID,
		AmountAvailable: req.AmountAvailable,
	}

	inventory, err := server.store.CreateInventory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventory)
}

type getCompanyProductInventoryRequest struct {
	CompanyID int64 `uri:"company_id" binding:"required,min=1"`
	ProductID int64 `uri:"product_id" binding:"required,min=1"`
}

func (server *Server) getCompanyProductInventory(ctx *gin.Context) {
	var req getCompanyProductInventoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetCompanyProductInventoryParams{
		CompanyID: req.CompanyID,
		ProductID: req.ProductID,
	}

	inventory, err := server.store.GetCompanyProductInventory(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventory)
}

type listCompanyInventoriesRequest struct {
	CompanyID int64 `form:"company_id" binding:"required,min=1"`
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCompanyInventories(ctx *gin.Context) {
	var req listCompanyInventoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCompanyInventoriesParams{
		CompanyID: req.CompanyID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	inventories, err := server.store.ListCompanyInventories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inventories)
}
