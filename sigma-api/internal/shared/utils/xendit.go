package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sigma-api/internal/shared/config"
)

type XenditInvoiceRequest struct {
	ExternalID      string  `json:"external_id"`
	Amount          float64 `json:"amount"`
	PayerEmail      string  `json:"payer_email"`
	Description     string  `json:"description"`
	InvoiceDuration int     `json:"invoice_duration"`
}

type XenditInvoiceResponse struct {
	ID         string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
	Status     string `json:"status"`
}

// CreateXenditInvoice calls the Xendit API to generate a payment link
func CreateXenditInvoice(req XenditInvoiceRequest) (*XenditInvoiceResponse, error) {
	url := "https://api.xendit.co/v2/invoices"
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Basic Auth: SecretKey as Username, empty Password
	auth := base64.StdEncoding.EncodeToString([]byte(config.AppConfig.XenditSecretKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("xendit api error: status %d", resp.StatusCode)
	}

	var xenditResp XenditInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&xenditResp); err != nil {
		return nil, err
	}

	return &xenditResp, nil
}
