package integration_test

import (
	"fmt"
	"go-loan-api/database"
	"go-loan-api/internal/client"
	"go-loan-api/internal/handler"
	"go-loan-api/internal/repository"
	"go-loan-api/internal/service"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testRouter *gin.Engine
var testDB *gorm.DB
var mockClient *MockCustomerClient

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Fatalf("cannot load .env.test: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	database.NewMigrate(db).Run()

	testDB = db
	mockClient = &MockCustomerClient{}
	testRouter = buildRouter(db, mockClient)

	os.Exit(m.Run())
}

func buildRouter(db *gorm.DB, customerClient client.CustomerClient) *gin.Engine {
	engine := gin.New()
	loanRepo := repository.NewLoanRepository(db)
	loanRuleRepo := repository.NewLoanRuleRepository(db)
	loanSvc := service.NewLoanService(loanRepo, loanRuleRepo, customerClient)
	loanHdl := handler.NewLoanHandler(loanSvc)
	loanHdl.RegisterRoute(engine)
	return engine
}

func TruncateAll(t *testing.T, tables ...string) {
	t.Helper()
	for _, table := range tables {
		if err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
			t.Fatalf("failed to truncate %s: %v", table, err)
		}
	}
}
