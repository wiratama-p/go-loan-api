package repository

import (
	"go-loan-api/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoanRuleRepository struct {
	DB *gorm.DB
}

func NewLoanRuleRepository(db *gorm.DB) *LoanRuleRepository {
	return &LoanRuleRepository{
		DB: db,
	}
}

func (repository *LoanRuleRepository) GetByProposedAmount(ctx *gin.Context, amount int64) *model.LoanRule {
	var loanRule model.LoanRule
	result := repository.DB.WithContext(ctx).
		Debug().
		Where("min_amount >= ?", amount).
		Order("max_amount ASC").
		Limit(1).
		Find(&loanRule)
	if result.RowsAffected == 0 {
		return nil
	}
	return &loanRule
}
