package zooz

// Result represents status and category of some methods response.
type Result struct {
	Status      ResultStatus `json:"status"`
	Category    string       `json:"category"`
	SubCategory string       `json:"sub_category"`
	Description string       `json:"description"`
}

type ResultStatus string
type TransactionCardEntryMode string

const (
	StatusSucceed ResultStatus = "Succeed"
	StatusFailed  ResultStatus = "Failed"

	// The initial transaction in which the customer agrees to using stored card information for subsequent customer-initiated transactions, or subsequent unscheduled transactions initiated by the merchant.
	ConsentTransaction TransactionCardEntryMode = "consent_transaction"
	// The initial transaction in which the customer agrees to using stored card information for subsequent scheduled (recurring) transactions.
	RecurringConsentTransaction TransactionCardEntryMode = "recurring_consent_transaction"
	// A transaction in a series of transactions that use stored card information and that are processed at fixed, regular intervals.
	RecurringSubsequentTransaction TransactionCardEntryMode = "recurring_subsequent_transaction"
	// Used for card-on-file transactions, initiated by the customer.
	CofCardholderInitiatedTransaction TransactionCardEntryMode = "cof_cardholder_initiated_transaction"
	// Used for unscheduled card-on-file transactions, initiated by you (as the merchant).
	CofMerchantInitiatedTransaction TransactionCardEntryMode = "cof_merchant_initiated_transaction"
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
	Internal *ThreeDSInternal `json:"internal,omitempty"`
	External *ThreeDSExternal `json:"external,omitempty"`
}

// ThreeDSInternal is a set of attributes for 3D-Secure 2.x.x handled internally by PaymentsOS
type ThreeDSInternal struct {
	DeviceChannel                    string `json:"device_channel,omitempty"`
	WorkPhone                        string `json:"work_phone,omitempty"`
	MobilePhone                      string `json:"mobile_phone,omitempty"`
	HomePhone                        string `json:"home_phone,omitempty"`
	MobilePhoneCountry               string `json:"mobile_phone_country,omitempty"`
	HomePhoneCountry                 string `json:"home_phone_country,omitempty"`
	WorkPhoneCountry                 string `json:"work_phone_country,omitempty"`
	AddressMatch                     *bool  `json:"address_match,omitempty"`
	ProductCode                      string `json:"product_code,omitempty"`
	ShippingMethodIndicator          string `json:"shipping_method_indicator,omitempty"`
	DeliveryTimeFrame                string `json:"delivery_time_frame,omitempty"`
	ReorderIndicator                 string `json:"reorder_indicator,omitempty"`
	PreOrderIndicator                string `json:"pre_order_indicator,omitempty"`
	PreOrderDate                     string `json:"pre_order_date,omitempty"`
	AccountAgeIndicator              string `json:"account_age_indicator,omitempty"`
	AccountCreateDate                string `json:"account_create_date,omitempty"`
	AccountChangeIndicator           string `json:"account_change_indicator,omitempty"`
	AccountChangeDate                string `json:"account_change_date,omitempty"`
	AccountPwdChangeIndicator        string `json:"account_pwd_change_indicator,omitempty"`
	AccountPwdChangeDate             string `json:"account_pwd_change_date,omitempty"`
	AccountAdditionalInformation     string `json:"account_additional_information,omitempty"`
	ShippingAddressUsageIndicator    string `json:"shipping_address_usage_indicator,omitempty"`
	ShippingAddressUsageDate         string `json:"shipping_address_usage_date,omitempty"`
	TransactionCountDay              string `json:"transaction_count_day,omitempty"`
	TransactionCountYear             string `json:"transaction_count_year,omitempty"`
	AddCardAttemptsDay               string `json:"add_card_attempts_day,omitempty"`
	AccountPurchasesSixMonths        string `json:"account_purchases_six_months,omitempty"`
	FraudActivity                    string `json:"fraud_activity,omitempty"`
	ShippingNameIndicator            string `json:"shipping_name_indicator,omitempty"`
	PaymentAccountIndicator          string `json:"payment_account_indicator,omitempty"`
	PaymentAccountAge                string `json:"payment_account_age,omitempty"`
	RequestorAuthenticationMethod    string `json:"requestor_authentication_method,omitempty"`
	RequestorAuthenticationTimestamp string `json:"requestor_authentication_timestamp,omitempty"`
	RequestorAuthenticationData      string `json:"requestor_authentication_data,omitempty"`
	PriorAuthenticationData          string `json:"prior_authentication_data,omitempty"`
	PriorAuthenticationMethod        string `json:"prior_authentication_method,omitempty"`
	PriorAuthenticationTimestamp     string `json:"prior_authentication_timestamp,omitempty"`
	PriorAuthenticationRef           string `json:"prior_authentication_ref,omitempty"`
	PurchaseDateTime                 string `json:"purchase_date_time,omitempty"`
	RecurringEndDate                 string `json:"recurring_end_date,omitempty"`
	RecurringFrequency               int64  `json:"recurring_frequency,omitempty"`
	BrowserHeader                    string `json:"browser_header,omitempty"`
	BrowserJavaEnabled               *bool  `json:"browser_java_enabled,omitempty"`
	BrowserLanguage                  string `json:"browser_language,omitempty"`
	BrowserColorDepth                string `json:"browser_color_depth,omitempty"`
	BrowserScreenHeight              string `json:"browser_screen_height,omitempty"`
	BrowserScreenWidth               string `json:"browser_screen_width,omitempty"`
	BrowserTimeZone                  string `json:"browser_time_zone,omitempty"`
	ChallengeIndicator               string `json:"challenge_indicator,omitempty"`
	ChallengeWindowSize              string `json:"challenge_window_size,omitempty"`
	SdkAppID                         string `json:"sdk_app_id,omitempty"`
	SdkEncryptedData                 string `json:"sdk_encrypted_data,omitempty"`
	SdkMaxTimeout                    string `json:"sdk_max_timeout,omitempty"`
	SdkReferenceNumber               string `json:"sdk_reference_number,omitempty"`
	SdkTransactionID                 string `json:"sdk_transaction_id,omitempty"`
	SdkInterface                     string `json:"sdk_interface,omitempty"`
	SdkUiType                        string `json:"sdk_ui_type,omitempty"`
	SdkEphemeralPublicKey            string `json:"sdk_ephemeral_public_key,omitempty"`
}

// ThreeDSInternal is a set of attributes for 3D-Secure 1.x.x and 2.x.x
type ThreeDSExternal struct {
	Version              string `json:"three_d_secure_version,omitempty"`
	AuthenticationStatus string `json:"three_d_secure_authentication_status,omitempty"`
	Encoding             string `json:"encoding"`
	XID                  string `json:"xid"`
	DSXID                string `json:"ds_xid,omitempty"`
	CAVV                 string `json:"cavv"`
	EciFlag              string `json:"eci_flag"`
}

// CofTransactionIndicators contains indicators pertaining to the use of a customer's stored card information
type CofTransactionIndicators struct {
	CardEntryMode           *TransactionCardEntryMode `json:"card_entry_mode,omitempty"`
	CofConsentTransactionId string                    `json:"cof_consent_transaction_id,omitempty"`
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
	CvvVerificationCode   string             `json:"cvv_verification_code"`
	TransactionCost       TransactionCost    `json:"transaction_cost"`
	AdditionalInformation map[string]string  `json:"additional_information"`
}

type TransactionCost struct {
	AppliedFees []AppliedFees `json:"applied_fees"`
}

type AppliedFees struct {
	Type         string  `json:"type"`
	Amount       int64   `json:"amount"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate"`
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
	Href string `json:"href"`
	*PaymentMethod
}
