package integration_test

import (
	"context"
	"go-loan-api/internal/client"
)

type MockCustomerClient struct {
	Customer *client.CustomerInfo
	Err      error
}

func (client *MockCustomerClient) GetCustomer(_ context.Context, _ string) (*client.CustomerInfo, error) {
	return client.Customer, client.Err
}
