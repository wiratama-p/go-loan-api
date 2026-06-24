package handler

import (
	"go-loan-api/internal/dto"
	"go-loan-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanHandler struct {
	loanService *service.LoanService
}

func NewLoanHandler(loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

func (handler *LoanHandler) RegisterRoute(engine *gin.Engine) {
	loansGroup := engine.Group("/loans")
	loansGroup.POST("", handler.Create)
	customersGroup := engine.Group("/customers")
	customersGroup.GET("/:id/loans", handler.GetByCustomerID)
}

func (handler *LoanHandler) Create(ctx *gin.Context) {
	var request dto.CreateLoanRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, &dto.Response[any]{
			Code:   http.StatusBadRequest,
			Status: "Validation Failed",
			Error:  err.Error(),
		})
		return
	}
	createdLoan, err := handler.loanService.Create(ctx, &request)
	if err != nil {
		ToErrorResponse(ctx, err)
		return
	}
	dto.Created(ctx, createdLoan)
}

func (handler *LoanHandler) GetByCustomerID(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	customerId := ctx.Param("id")

	loans := handler.loanService.GetAllByCustomerID(ctx, page, customerId)
	dto.Ok(ctx, loans)
}
