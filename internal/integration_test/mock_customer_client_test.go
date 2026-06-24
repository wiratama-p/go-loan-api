package integration_test

import (
	"context"
	"go-loan-api/internal/client_response"
)

type MockCustomerClient struct {
	Customer *client_response.CustomerInfo
	Err      error
}

func (m *MockCustomerClient) GetCustomer(_ context.Context, _ string) (*client_response.CustomerInfo, error) {
	return m.Customer, m.Err
}
