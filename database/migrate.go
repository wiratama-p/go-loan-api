package database

import (
	"go-loan-api/internal/model"
	"log"

	"gorm.io/gorm"
)

type Migrate struct {
	DB *gorm.DB
}

func NewMigrate(db *gorm.DB) *Migrate {
	return &Migrate{DB: db}
}

func (migrate *Migrate) Run() {
	err := migrate.DB.AutoMigrate(model.Loan{}, model.LoanRule{})
	if err != nil {
		log.Panic("error when auto migrate DB", err.Error())
		return
	}
}
