package repository

import (
	"go-loan-api/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoanRepository struct {
	DB *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{
		DB: db,
	}
}

func (repository LoanRepository) Create(ctx *gin.Context, loan *model.Loan) (*model.Loan, error) {
	result := repository.DB.WithContext(ctx).Create(&loan)
	return loan, result.Error
}

func (repository *LoanRepository) GetAllByCustomerId(ctx *gin.Context, page int, customerID string) *[]model.Loan {
	var loans []model.Loan
	offset := (page - 1) * 10
	repository.DB.
		Debug().
		WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("COALESCE(updated_at, created_at) DESC").
		Offset(offset).
		Limit(10).
		Find(&loans)

	return &loans
}
