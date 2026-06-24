package dto

import (
	"go-loan-api/internal/client_response"
	"go-loan-api/internal/model"

	"github.com/google/uuid"
)

type CreateLoanRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
	Purpose    string `json:"purpose"`
}

type UpdateLoanStatusRequest struct {
	Status string `json:"status" binding:"required" oneof:"PENDING APPROVED REJECTED"`
}

type LoanResponse struct {
	ID           string  `json:"id"`
	CustomerID   *string `json:"customer_id"`
	CustomerName *string `json:"customer_name"`
	Amount       int64   `json:"amount"`
	InterestRate int     `json:"interest_rate"`
	Tenure       int     `json:"tenure"`
	Status       string  `json:"status"`
	Purpose      string  `json:"purpose"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
}

func ToLoanModel(request *CreateLoanRequest, rule *model.LoanRule) *model.Loan {
	customerId, _ := uuid.Parse(request.CustomerID)
	return &model.Loan{
		CustomerID:   customerId,
		Amount:       request.Amount,
		InterestRate: rule.InterestRate,
		Tenure:       rule.Tenure,
		Purpose:      request.Purpose,
		Status:       "PENDING",
	}
}

func ToLoanResponse(loan *model.Loan, customerInfo *client_response.CustomerInfo) *LoanResponse {
	return &LoanResponse{
		ID:           loan.ID.String(),
		CustomerID:   customerInfo.ID,
		CustomerName: customerInfo.Name,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Tenure:       loan.Tenure,
		Status:       loan.Status,
		Purpose:      loan.Purpose,
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
}

func ToLoanResponses(loans *[]model.Loan) *[]LoanResponse {
	responses := []LoanResponse{}
	for _, loan := range *loans {
		loanResponse := ToLoanResponse(&loan, &client_response.CustomerInfo{
			ID:   nil,
			Name: nil,
		})
		responses = append(responses, *loanResponse)
	}
	return &responses
}
