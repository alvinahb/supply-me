package api

import (
	"database/sql"
	"net/http"

	db "github.com/alvinahb/supply-me/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createCompanyRequest struct {
	CompanyType string `json:"company_type" binding:"required,oneof=Restaurant Fournisseur"`
	CompanyName string `json:"company_name" binding:"required"`
}

func (server *Server) createCompany(ctx *gin.Context) {
	var req createCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCompanyParams{
		CompanyType: req.CompanyType,
		CompanyName: req.CompanyName,
	}

	company, err := server.store.CreateCompany(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

type getCompanyRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCompany(ctx *gin.Context) {
	var req getCompanyRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	company, err := server.store.GetCompany(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

type listCompaniesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCompanies(ctx *gin.Context) {
	var req listCompaniesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCompaniesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	companies, err := server.store.ListCompanies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, companies)
}
