package model

import "github.com/google/uuid"

type Loan struct {
	ID            uuid.UUID `gorm:"type:uuid; default:gen_random_uuid(); primaryKey"`
	CustomerID    uuid.UUID `gorm:"type:uuid; not null; index"`
	Amount        int64     `gorm:"type:bigint; not null"`
	InterestRate  int       `gorm:"type:int; not null"`
	Tenure        int       `gorm:"type:int; not null"`
	Status        string    `gorm:"type:varchar(50); not null; default:'PENDING'"`
	PaymentStatus string    `gorm:"type:varchar(50); not null; default:'OUTSTANDING'"`
	Purpose       string
	CreatedAt     int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt     int64 `gorm:"autoUpdateTime:milli"`
}

func (Loan) TableName() string {
	return "loan"
}
