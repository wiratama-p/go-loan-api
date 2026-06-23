package handler

import (
	"go-loan-api/internal/dto"
	"go-loan-api/internal/service"
	"net/http"

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
	group := engine.Group("/loans")
	group.POST("", handler.Create)
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
