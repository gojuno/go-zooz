package zooz

import (
	"context"
	"net/http"
)

const supportedPaymentMethodsPath = "supported-payment-methods"

type SupportedPaymentMethodClient struct {
	Caller Caller
}

type SupportedPaymentMethod struct {
	ConfigurationID         string              `json:"configuration_id"`
	ConfigurationName       string              `json:"configuration_name"`
	ProviderID              string              `json:"provider_id"`
	ProviderName            string              `json:"provider_name"`
	SupportedPaymentMethods []PaymentMethodData `json:"supported_payment_methods"`
	Result                  Result              `json:"result"`
}

type PaymentMethodData struct {
	DisplayName string                  `json:"display_name"`
	Vendor      string                  `json:"vendor"`
	SourceType  SourceType              `json:"source_type"`
	Status      PaymentMethodDataStatus `json:"status"`
	Country     string                  `json:"country"`
	LogoURL     string                  `json:"logo_url"`
	Amounts     []Amount                `json:"amounts"`
}

type SourceType string

const (
	BankTransfer SourceType = "bank_transfer"
	Cash         SourceType = "cash"
	Ewallet      SourceType = "ewallet"
	DebitDirect  SourceType = "debit_redirect"
	Loyalty      SourceType = "loyalty"
	CreditCard   SourceType = "credit_card"
	Credit       SourceType = "credit"
)

type PaymentMethodDataStatus string

const (
	PMStatusAvailable            PaymentMethodDataStatus = "available"
	PMStatusTemporaryUnavailable PaymentMethodDataStatus = "temporarily_unavailable"
	PMStatusDisabled             PaymentMethodDataStatus = "disabled"
)

type Amount struct {
	Min      int64  `json:"min"`
	Max      int64  `json:"max"`
	Currency string `json:"currency"`
}

func (c *SupportedPaymentMethodClient) Get(ctx context.Context) ([]SupportedPaymentMethod, error) {
	var spms []SupportedPaymentMethod
	err := c.Caller.Call(ctx, http.MethodGet, supportedPaymentMethodsPath, nil, nil, &spms)
	if err != nil {
		return nil, err
	}
	return spms, nil
}
