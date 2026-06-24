package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-loan-api/internal/client_response"
	"go-loan-api/internal/dto"
	"go-loan-api/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoanHandler(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		t.Run("When loan rule exists and customer exists, should create loan successfully", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")
			customerID := "550e8400-e29b-41d4-a716-446655440000"
			customerName := "Alice Doe"
			mockClient.Customer = &client_response.CustomerInfo{ID: &customerID, Name: &customerName}
			mockClient.Err = nil
			require.NoError(t, testDB.Create(&model.LoanRule{MinAmount: 1000, MaxAmount: 100000, InterestRate: 5, Tenure: 12}).Error)

			body := dto.CreateLoanRequest{CustomerID: customerID, Amount: 500, Purpose: "Business"}

			result := post(t, "/loans", body)

			var response dto.Response[dto.LoanResponse]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusCreated, result.Code)
			assert.Equal(t, customerID, *response.Data.CustomerID)
			assert.Equal(t, int64(500), response.Data.Amount)
			assert.Equal(t, "PENDING", response.Data.Status)
			assert.Equal(t, 5, response.Data.InterestRate)
			assert.Equal(t, 12, response.Data.Tenure)

			var loan model.Loan
			require.NoError(t, testDB.First(&loan, "id = ?", response.Data.ID).Error)
			assert.Equal(t, int64(500), loan.Amount)
		})

		t.Run("When required field is missing, should return bad request", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")

			result := post(t, "/loans", map[string]any{"purpose": "Business"})

			var response dto.Response[any]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusBadRequest, result.Code)
			assert.NotEmpty(t, response.Error)
		})

		t.Run("When no loan rule matches the amount, should return bad request", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")

			body := dto.CreateLoanRequest{CustomerID: "550e8400-e29b-41d4-a716-446655440000", Amount: 50000}

			result := post(t, "/loans", body)

			var response dto.Response[any]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusBadRequest, result.Code)
			assert.Contains(t, response.Error, "not found")
		})

		t.Run("When customer service returns error, should return internal server error", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")
			mockClient.Customer = nil
			mockClient.Err = errors.New("customer service unavailable")
			require.NoError(t, testDB.Create(&model.LoanRule{MinAmount: 1000, MaxAmount: 100000, InterestRate: 5, Tenure: 12}).Error)

			body := dto.CreateLoanRequest{CustomerID: "550e8400-e29b-41d4-a716-446655440000", Amount: 500}

			result := post(t, "/loans", body)

			assert.Equal(t, http.StatusInternalServerError, result.Code)
		})
	})

	t.Run("GetByCustomerID", func(t *testing.T) {
		t.Run("When customer has loans, should return list of loans", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")
			customerID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")

			loans := []model.Loan{
				{CustomerID: customerID, Amount: 5000, InterestRate: 5, Tenure: 12, Status: "PENDING", PaymentStatus: "OUTSTANDING"},
				{CustomerID: customerID, Amount: 10000, InterestRate: 8, Tenure: 24, Status: "APPROVED", PaymentStatus: "OUTSTANDING"},
			}
			for i := range loans {
				require.NoError(t, testDB.Create(&loans[i]).Error)
			}

			result := get(t, "/customers/"+customerID.String()+"/loans?page=1")

			var response dto.Response[[]dto.LoanResponse]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusOK, result.Code)
			assert.Len(t, *response.Data, 2)
		})

		t.Run("When customer has no loans, should return empty list", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")

			result := get(t, "/customers/non-existent-id/loans?page=1")

			var response dto.Response[[]dto.LoanResponse]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusOK, result.Code)
			assert.Empty(t, *response.Data)
		})
	})

	t.Run("UpdateLoanStatus", func(t *testing.T) {
		t.Run("When loan exists, should update loan status successfully", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")
			customerID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
			loan := model.Loan{
				CustomerID:    customerID,
				Amount:        5000,
				InterestRate:  5,
				Tenure:        12,
				Status:        "PENDING",
				PaymentStatus: "OUTSTANDING",
			}
			require.NoError(t, testDB.Create(&loan).Error)

			result := patch(t, "/loans/"+loan.ID.String()+"/status", dto.UpdateLoanStatusRequest{Status: "APPROVED"})

			var response dto.Response[dto.LoanResponse]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusOK, result.Code)
			assert.Equal(t, "APPROVED", response.Data.Status)

			var updatedLoan model.Loan
			require.NoError(t, testDB.First(&updatedLoan, "id = ?", loan.ID).Error)
			assert.Equal(t, "APPROVED", updatedLoan.Status)
		})

		t.Run("When loan does not exist, should return not found", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")

			result := patch(t, "/loans/non-existent-id/status", dto.UpdateLoanStatusRequest{Status: "APPROVED"})

			var response dto.Response[any]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusNotFound, result.Code)
		})

		t.Run("When status field is missing, should return bad request", func(t *testing.T) {
			TruncateAll(t, "loan", "loan_rule")

			result := patch(t, "/loans/some-id/status", map[string]any{})

			var response dto.Response[any]
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &response))
			assert.Equal(t, http.StatusBadRequest, result.Code)
			assert.NotEmpty(t, response.Error)
		})
	})
}

func post(t *testing.T, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func patch(t *testing.T, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}
