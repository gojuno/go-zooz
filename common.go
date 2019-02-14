package zooz

// Result represents status and category of some methods response.
type Result struct {
	Status      ResultStatus `json:"status"`
	Category    string       `json:"category"`
	SubCategory string       `json:"sub_category"`
	Description string       `json:"description"`
}

type ResultStatus string

const (
	StatusSucceed ResultStatus = "Succeed"
	StatusFailed  ResultStatus = "Failed"
)

// ClientInfo represents optional request params for some methods.
type ClientInfo struct {
	IPAddress string
	UserAgent string
}

// Address is a set of fields describing customer address.
type Address struct {
	Country   string `json:"country,omitempty"`
	State     string `json:"state,omitempty"`
	City      string `json:"city,omitempty"`
	Line1     string `json:"line1,omitempty"`
	Line2     string `json:"line2,omitempty"`
	ZipCode   string `json:"zip_code,omitempty"`
	Title     string `json:"title,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
}

// AdditionalDetails is a set of any custom key-value info.
type AdditionalDetails map[string]string

// ThreeDSecureAttributes is a set of attributes for 3D-Secure.
type ThreeDSecureAttributes struct {
	Encoding string `json:"encoding"`
	XID      string `json:"xid"`
	CAVV     string `json:"cavv"`
	EciFlag  string `json:"eci_flag"`
}

// Installments is a set of options of installments.
type Installments struct {
	NumberOfInstallments    int64 `json:"number_of_installments"`
	FirstPaymentAmount      int64 `json:"first_payment_amount"`
	RemainingPaymentsAmount int64 `json:"remaining_payments_amount"`
}

// ProviderData is a set of params describing payment provider.
type ProviderData struct {
	ProviderName          string             `json:"provider_name"`
	ResponseCode          string             `json:"response_code"`
	Description           string             `json:"description"`
	RawResponse           string             `json:"raw_response"`
	AvsCode               string             `json:"avs_code"`
	AuthorizationCode     string             `json:"authorization_code"`
	TransactionID         string             `json:"transaction_id"`
	ExternalID            string             `json:"external_id"`
	Documents             []ProviderDocument `json:"documents"`
	AdditionalInformation map[string]string  `json:"additional_information"`
}

// ProviderDocument represents provider document.
type ProviderDocument struct {
	Descriptor  string `json:"descriptor"`
	ContentType string `json:"content_type"`
	Href        string `json:"href"`
}

// PaymentMethodDetails represents payment method details for POST requests.
type PaymentMethodDetails struct {
	Type              string            `json:"type"`
	Token             string            `json:"token,omitempty"`
	CreditCardCvv     string            `json:"credit_card_cvv,omitempty"`
	SourceType        string            `json:"source_type,omitempty"`
	Vendor            string            `json:"vendor,omitempty"`
	AdditionalDetails AdditionalDetails `json:"additional_details,omitempty"`
}

// PaymentMethodHref wraps PaymentMethod with associated href.
type PaymentMethodHref struct {
	Href           string `json:"href"`
	*PaymentMethod
}
