package model

import "github.com/google/uuid"

type LoanRule struct {
	ID           uuid.UUID `gorm:"type:uuid; default:gen_random_uuid(); primaryKey"`
	MinAmount    int64     `gorm:"type:int; not null"`
	MaxAmount    int64     `gorm:"type:int; not null"`
	InterestRate int       `gorm:"type:int; not null"`
	Tenure       int       `gorm:"type:int; not null"`
	CreatedAt    int64     `gorm:"autoCreateTime:milli"`
	UpdatedAt    int64     `gorm:"autoUpdateTime:milli"`
}

func (LoanRule) TableName() string {
	return "loan_rule"
}
