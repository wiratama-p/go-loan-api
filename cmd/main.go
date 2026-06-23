package main

import (
	"go-loan-api/database"
	"go-loan-api/internal/client"
	"go-loan-api/internal/handler"
	"go-loan-api/internal/repository"
	"go-loan-api/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	routeEngine := gin.Default()

	DB := database.Connect()

	migrate := database.NewMigrate(DB)
	migrate.Run()

	customerClient := client.NewHttpCustomerClient()
	loanRepository := repository.NewLoanRepository(DB)
	loanRuleRepository := repository.NewLoanRuleRepository(DB)
	loanService := service.NewLoanService(loanRepository, loanRuleRepository, customerClient)
	customerHandler := handler.NewLoanHandler(loanService)
	customerHandler.RegisterRoute(routeEngine)

	err := routeEngine.Run(":8081")
	if err != nil {
		log.Panic("Error running server", err)
		return
	}
}
