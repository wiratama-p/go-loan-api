package service

import (
	"fmt"
	"go-loan-api/internal/apperror"
	"go-loan-api/internal/client"
	"go-loan-api/internal/dto"
	"go-loan-api/internal/repository"

	"github.com/gin-gonic/gin"
)

var (
	LoanRuleNotFound = "Loan rule for amount %d not found"
)

type LoanService struct {
	loanRepository     *repository.LoanRepository
	loanRuleRepository *repository.LoanRuleRepository
	customerClient     client.CustomerClient
}

func NewLoanService(loanRepository *repository.LoanRepository, loanRuleRepository *repository.LoanRuleRepository, customerClient client.CustomerClient) *LoanService {
	return &LoanService{
		loanRepository:     loanRepository,
		loanRuleRepository: loanRuleRepository,
		customerClient:     customerClient,
	}
}

func (loanService *LoanService) Create(ctx *gin.Context, request *dto.CreateLoanRequest) (*dto.LoanResponse, error) {
	loanRule := loanService.loanRuleRepository.GetByProposedAmount(ctx, request.Amount)
	if loanRule == nil {
		return nil, apperror.BadRequest(fmt.Sprintf(LoanRuleNotFound, request.Amount))
	}
	foundCustomer, err := loanService.customerClient.GetCustomer(ctx.Request.Context(), request.CustomerID)
	if err != nil {
		return nil, err
	}
	loan := dto.ToLoanModel(request, loanRule)
	createdLoan, _ := loanService.loanRepository.Create(ctx, loan)

	return dto.ToLoanResponse(createdLoan, foundCustomer), nil
}

func (loanService *LoanService) GetAllByCustomerID(ctx *gin.Context, page int, customerID string) *[]dto.LoanResponse {
	loans := loanService.loanRepository.GetAllByCustomerId(ctx, page, customerID)

	return dto.ToLoanResponses(loans)
}
