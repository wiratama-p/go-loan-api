package client

import (
	"context"
	"encoding/json"
	"fmt"
	"go-loan-api/internal/apperror"
	"go-loan-api/internal/dto"
	"net/http"
	"os"
	"time"
)

var (
	ErrCustomerNotFound       = "customer not found"
	ErrCustomerSvcUnavailable = "customer service unavailable"
)

type CustomerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CustomerClient interface {
	GetCustomer(ctx context.Context, customerID string) (*CustomerInfo, error)
}

type HttpCustomerClient struct {
	BaseUrl    string
	HttpClient *http.Client
}

func NewHttpCustomerClient() *HttpCustomerClient {
	timeout, _ := time.ParseDuration(os.Getenv("CLIENT_CUSTOMER_SVC_TIMEOUT"))
	return &HttpCustomerClient{
		BaseUrl: os.Getenv("CLIENT_CUSTOMER_SVC_BASE_URL"),
		HttpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (client *HttpCustomerClient) GetCustomer(ctx context.Context, customerID string) (*CustomerInfo, error) {
	url := fmt.Sprintf("%s/customers/%s", client.BaseUrl, customerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, apperror.InternalServerError(ErrCustomerSvcUnavailable)
	}
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, apperror.InternalServerError(ErrCustomerSvcUnavailable)
	}
	defer resp.Body.Close()

	if http.StatusNotFound == resp.StatusCode {
		return nil, apperror.InternalServerError(ErrCustomerNotFound)
	}
	if http.StatusOK == resp.StatusCode {
		var customer dto.Response[CustomerInfo]
		if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
			return nil, apperror.InternalServerError(ErrCustomerSvcUnavailable)
		}
		return customer.Data, nil
	}
	return nil, apperror.InternalServerError(ErrCustomerSvcUnavailable)
}
