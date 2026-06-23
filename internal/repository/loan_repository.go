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
