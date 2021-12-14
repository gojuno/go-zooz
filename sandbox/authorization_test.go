package sandbox

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorization(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new (successful, all fields, no 3ds, no customer) & get", func(t *testing.T) {
		t.Parallel()

		token, tokenParams := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		reconciliationID := randomString(32)

		authorizationCreated, err := client.Authorization().New(
			context.Background(),
			randomString(32),
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token.Token,
					CreditCardCvv: token.EncryptedCVV,
				},
				MerchantSiteURL:        "abc",
				ReconciliationID:       reconciliationID,
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderSpecificData: map[string]interface{}{
					"provider data 1": "aaa",
					"provider data 2": 123,
					"provider data 3": map[string]interface{}{
						"provider data 4": "bbb",
					},
				},
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				COFTransactionIndicators: &zooz.COFTransactionIndicators{
					CardEntryMode:           "consent_transaction",
					COFConsentTransactionID: "",
				},
			},
			&zooz.ClientInfo{
				IPAddress: "95.173.136.70",
				UserAgent: "my-user-agent",
			},
		)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, authorizationCreated.ID)
			assert.NotEmpty(t, authorizationCreated.Created)
			assert.NotEmpty(t, authorizationCreated.PaymentMethod.FingerPrint)
			assert.NotEmpty(t, authorizationCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, authorizationCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, authorizationCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, authorizationCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, authorizationCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Authorization{
				ID: authorizationCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Amount:           5000,
				Created:          authorizationCreated.Created, // ignore
				ReconciliationID: reconciliationID,
				PaymentMethod: zooz.PaymentMethod{
					Href:               "", // empty, because is not assigned to customer
					Type:               "tokenized",
					TokenType:          string(token.TokenType),
					PassLuhnValidation: token.PassLuhnValidation,
					Token:              token.Token,
					Created:            json.Number(token.Created),
					Customer:           "",
					AdditionalDetails:  nil, // why empty? is not empty on get
					BinNumber:          json.Number(token.BinNumber),
					Vendor:             token.Vendor,
					Issuer:             token.Issuer,
					CardType:           token.CardType,
					Level:              token.Level,
					CountryCode:        token.CountryCode,
					HolderName:         tokenParams.HolderName,
					ExpirationDate:     normalizeExpirationDate(tokenParams.ExpirationDate),
					Last4Digits:        last4(tokenParams.CardNumber),
					ShippingAddress:    nil, // why empty?
					BillingAddress: &zooz.Address{
						Country:   tokenParams.BillingAddress.Country,
						State:     tokenParams.BillingAddress.State,
						City:      tokenParams.BillingAddress.City,
						Line1:     tokenParams.BillingAddress.Line1,
						Line2:     tokenParams.BillingAddress.Line2,
						ZipCode:   tokenParams.BillingAddress.ZipCode,
						Title:     "", // why empty?
						FirstName: tokenParams.BillingAddress.FirstName,
						LastName:  tokenParams.BillingAddress.LastName,
						Phone:     tokenParams.BillingAddress.Phone,
						Email:     tokenParams.BillingAddress.Email,
					},
					FingerPrint: authorizationCreated.PaymentMethod.FingerPrint, // ignore
				},
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Authorized.",
					RawResponse:           authorizationCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     authorizationCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         authorizationCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            authorizationCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				ProviderSpecificData:       zooz.DecodedJSON{},                         // why empty?
				ProviderConfiguration:      authorizationCreated.ProviderConfiguration, // ignore
				OriginatingPurchaseCountry: "RUS",                                      // based on IP
				IPAddress:                  "95.173.136.70",
				Redirection:                nil,
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				DecisionEngineExecution: authorizationCreated.DecisionEngineExecution, // ignore, since it depends on the provider list
			}, authorizationCreated)
		})

		authorizationRetrieved, err := client.Authorization().Get(context.Background(), payment.ID, authorizationCreated.ID)
		require.NoError(t, err)
		must(t, func() {
			// empty on New, correct on Get
			assert.Equal(t, tokenParams.AdditionalDetails, authorizationRetrieved.PaymentMethod.AdditionalDetails)
			authorizationRetrieved.PaymentMethod.AdditionalDetails = nil

			// empty on New, nil on Get
			assert.Nil(t, authorizationRetrieved.ProviderSpecificData)
			authorizationRetrieved.ProviderSpecificData = zooz.DecodedJSON{}

			// decision engine exectuion is not returned when retreving auth by id
			authorizationCreated.DecisionEngineExecution = zooz.DecisionEngineExecution{}

			assert.Equal(t, authorizationCreated, authorizationRetrieved)
		})
	})

	t.Run("new (successful, all fields, no 3ds, with customer, token not assigned) & get", func(t *testing.T) {
		t.Parallel()

		token, tokenParams := PrepareToken(t, client)
		customer := PrepareCustomer(t, client)
		payment := PreparePayment(t, client, 5000, customer)
		reconciliationID := randomString(32)

		authorizationCreated, err := client.Authorization().New(
			context.Background(),
			randomString(32),
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token.Token,
					CreditCardCvv: token.EncryptedCVV,
				},
				MerchantSiteURL:        "abc",
				ReconciliationID:       reconciliationID,
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderSpecificData: map[string]interface{}{
					"provider data 1": "aaa",
					"provider data 2": 123,
					"provider data 3": map[string]interface{}{
						"provider data 4": "bbb",
					},
				},
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				COFTransactionIndicators: &zooz.COFTransactionIndicators{
					CardEntryMode:           "consent_transaction",
					COFConsentTransactionID: "",
				},
			},
			&zooz.ClientInfo{
				IPAddress: "95.173.136.70",
				UserAgent: "my-user-agent",
			},
		)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, authorizationCreated.ID)
			assert.NotEmpty(t, authorizationCreated.Created)
			assert.NotEmpty(t, authorizationCreated.PaymentMethod.FingerPrint)
			assert.NotEmpty(t, authorizationCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, authorizationCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, authorizationCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, authorizationCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, authorizationCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Authorization{
				ID: authorizationCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Amount:           5000,
				Created:          authorizationCreated.Created, // ignore
				ReconciliationID: reconciliationID,
				PaymentMethod: zooz.PaymentMethod{
					Href:               "", // empty, because is not assigned to customer
					Type:               "tokenized",
					TokenType:          string(token.TokenType),
					PassLuhnValidation: token.PassLuhnValidation,
					Token:              token.Token,
					Created:            json.Number(token.Created),
					Customer:           "",
					AdditionalDetails:  nil, // why empty? is not empty on get
					BinNumber:          json.Number(token.BinNumber),
					Vendor:             token.Vendor,
					Issuer:             token.Issuer,
					CardType:           token.CardType,
					Level:              token.Level,
					CountryCode:        token.CountryCode,
					HolderName:         tokenParams.HolderName,
					ExpirationDate:     normalizeExpirationDate(tokenParams.ExpirationDate),
					Last4Digits:        last4(tokenParams.CardNumber),
					ShippingAddress:    nil, // why empty?
					BillingAddress: &zooz.Address{
						Country:   tokenParams.BillingAddress.Country,
						State:     tokenParams.BillingAddress.State,
						City:      tokenParams.BillingAddress.City,
						Line1:     tokenParams.BillingAddress.Line1,
						Line2:     tokenParams.BillingAddress.Line2,
						ZipCode:   tokenParams.BillingAddress.ZipCode,
						Title:     "", // why empty?
						FirstName: tokenParams.BillingAddress.FirstName,
						LastName:  tokenParams.BillingAddress.LastName,
						Phone:     tokenParams.BillingAddress.Phone,
						Email:     tokenParams.BillingAddress.Email,
					},
					FingerPrint: authorizationCreated.PaymentMethod.FingerPrint, // ignore
				},
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Authorized.",
					RawResponse:           authorizationCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     authorizationCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         authorizationCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            authorizationCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				ProviderSpecificData:       zooz.DecodedJSON{},                         // why empty?
				ProviderConfiguration:      authorizationCreated.ProviderConfiguration, // ignore
				OriginatingPurchaseCountry: "RUS",                                      // based on IP
				IPAddress:                  "95.173.136.70",
				Redirection:                nil,
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				DecisionEngineExecution: authorizationCreated.DecisionEngineExecution, // ignore
			}, authorizationCreated)
		})

		authorizationRetrieved, err := client.Authorization().Get(context.Background(), payment.ID, authorizationCreated.ID)
		require.NoError(t, err)
		must(t, func() {
			// empty on New, correct on Get
			assert.Equal(t, tokenParams.AdditionalDetails, authorizationRetrieved.PaymentMethod.AdditionalDetails)
			authorizationRetrieved.PaymentMethod.AdditionalDetails = nil

			// empty on New, nil on Get
			assert.Nil(t, authorizationRetrieved.ProviderSpecificData)
			authorizationRetrieved.ProviderSpecificData = zooz.DecodedJSON{}

			// decision engine exectuion is not returned when retreving auth by id
			authorizationCreated.DecisionEngineExecution = zooz.DecisionEngineExecution{}

			assert.Equal(t, authorizationCreated, authorizationRetrieved)
		})
	})

	t.Run("new (successful, all fields, no 3ds, token assigned to customer) & get", func(t *testing.T) {
		t.Parallel()

		token, tokenParams := PrepareToken(t, client)
		customer := PrepareCustomer(t, client)
		paymentMethod, err := client.PaymentMethod().New(context.Background(), randomString(32), customer.ID, token.Token)
		require.NoError(t, err)
		payment := PreparePayment(t, client, 5000, customer)
		reconciliationID := randomString(32)

		authorizationCreated, err := client.Authorization().New(
			context.Background(),
			randomString(32),
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token.Token,
					CreditCardCvv: token.EncryptedCVV,
				},
				MerchantSiteURL:        "abc",
				ReconciliationID:       reconciliationID,
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderSpecificData: map[string]interface{}{
					"provider data 1": "aaa",
					"provider data 2": 123,
					"provider data 3": map[string]interface{}{
						"provider data 4": "bbb",
					},
				},
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				COFTransactionIndicators: &zooz.COFTransactionIndicators{
					CardEntryMode:           "consent_transaction",
					COFConsentTransactionID: "",
				},
			},
			&zooz.ClientInfo{
				IPAddress: "95.173.136.70",
				UserAgent: "my-user-agent",
			},
		)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, authorizationCreated.ID)
			assert.NotEmpty(t, authorizationCreated.Created)
			assert.NotEmpty(t, authorizationCreated.PaymentMethod.FingerPrint)
			assert.NotEmpty(t, authorizationCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, authorizationCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, authorizationCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, authorizationCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, authorizationCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Authorization{
				ID: authorizationCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Amount:           5000,
				Created:          authorizationCreated.Created, // ignore
				ReconciliationID: reconciliationID,
				PaymentMethod: zooz.PaymentMethod{
					Href:               "", // why empty?
					Type:               paymentMethod.Type,
					TokenType:          paymentMethod.TokenType,
					PassLuhnValidation: paymentMethod.PassLuhnValidation,
					Token:              paymentMethod.Token,
					Created:            paymentMethod.Created,
					Customer:           "",  // why empty?
					AdditionalDetails:  nil, // why empty? is not empty on get
					BinNumber:          paymentMethod.BinNumber,
					Vendor:             paymentMethod.Vendor,
					Issuer:             paymentMethod.Issuer,
					CardType:           paymentMethod.CardType,
					Level:              paymentMethod.Level,
					CountryCode:        paymentMethod.CountryCode,
					HolderName:         paymentMethod.HolderName,
					ExpirationDate:     paymentMethod.ExpirationDate,
					Last4Digits:        paymentMethod.Last4Digits,
					ShippingAddress:    paymentMethod.ShippingAddress,
					BillingAddress: &zooz.Address{
						Country:   paymentMethod.BillingAddress.Country,
						State:     paymentMethod.BillingAddress.State,
						City:      paymentMethod.BillingAddress.City,
						Line1:     paymentMethod.BillingAddress.Line1,
						Line2:     paymentMethod.BillingAddress.Line2,
						ZipCode:   paymentMethod.BillingAddress.ZipCode,
						Title:     "", // why empty?
						FirstName: paymentMethod.BillingAddress.FirstName,
						LastName:  paymentMethod.BillingAddress.LastName,
						Phone:     paymentMethod.BillingAddress.Phone,
						Email:     paymentMethod.BillingAddress.Email,
					},
					FingerPrint: paymentMethod.FingerPrint,
				},
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Authorized.",
					RawResponse:           authorizationCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     authorizationCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         authorizationCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            authorizationCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				ProviderSpecificData:       zooz.DecodedJSON{},                         // why empty?
				ProviderConfiguration:      authorizationCreated.ProviderConfiguration, // ignore
				OriginatingPurchaseCountry: "RUS",                                      // based on IP
				IPAddress:                  "95.173.136.70",
				Redirection:                nil,
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				DecisionEngineExecution: authorizationCreated.DecisionEngineExecution,
			}, authorizationCreated)
		})

		authorizationRetrieved, err := client.Authorization().Get(context.Background(), payment.ID, authorizationCreated.ID)
		require.NoError(t, err)
		must(t, func() {
			// empty on New, correct on Get
			assert.Equal(t, tokenParams.AdditionalDetails, authorizationRetrieved.PaymentMethod.AdditionalDetails)
			authorizationRetrieved.PaymentMethod.AdditionalDetails = nil

			// empty on New, nil on Get
			assert.Nil(t, authorizationRetrieved.ProviderSpecificData)
			authorizationRetrieved.ProviderSpecificData = zooz.DecodedJSON{}

			// decision engine exectuion is not returned when retreving auth by id
			authorizationCreated.DecisionEngineExecution = zooz.DecisionEngineExecution{}

			assert.Equal(t, authorizationCreated, authorizationRetrieved)
		})
	})

	t.Run("new (successful, required fields)", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)

		authorizationCreated, err := client.Authorization().New(
			context.Background(),
			randomString(32),
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:  "tokenized",
					Token: token.Token,
				},
				MerchantSiteURL:          "",
				ReconciliationID:         "",
				ThreeDSecureAttributes:   nil,
				Installments:             nil,
				ProviderSpecificData:     nil,
				AdditionalDetails:        nil,
				COFTransactionIndicators: nil,
			},
			nil,
		)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, authorizationCreated.ID)
			assert.NotEmpty(t, authorizationCreated.Created)
			assert.NotEmpty(t, authorizationCreated.PaymentMethod.FingerPrint)
			assert.NotEmpty(t, authorizationCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Authorization{
				ID: authorizationCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Amount:                     5000,
				Created:                    authorizationCreated.Created, // ignore
				ReconciliationID:           "",
				PaymentMethod:              authorizationCreated.PaymentMethod, // ignore
				ThreeDSecureAttributes:     nil,
				Installments:               nil,
				ProviderData:               authorizationCreated.ProviderData, // ignore
				ProviderSpecificData:       zooz.DecodedJSON{},
				ProviderConfiguration:      authorizationCreated.ProviderConfiguration, // ignore
				OriginatingPurchaseCountry: "",                                         // empty because no IP
				IPAddress:                  "",
				Redirection:                nil,
				AdditionalDetails:          nil,
				DecisionEngineExecution:    authorizationCreated.DecisionEngineExecution,
			}, authorizationCreated)
		})
	})

	t.Run("new (failed, all fields, no 3ds, no customer) & get", func(t *testing.T) {
		t.Parallel()

		const cardNumber = "5555555555554444" // to fail authorization, see https://developers.paymentsos.com/docs/testing/mockprovider-reference.html

		token, err := client.CreditCardToken().New(context.Background(), randomString(32), &zooz.CreditCardTokenParams{
			HolderName:    "name",
			CardNumber:    cardNumber,
			CreditCardCVV: "123",
		})
		require.NoError(t, err)

		payment := PreparePayment(t, client, 5000, nil)
		reconciliationID := randomString(32)

		authorizationCreated, err := client.Authorization().New(
			context.Background(),
			randomString(32),
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token.Token,
					CreditCardCvv: token.EncryptedCVV,
				},
				MerchantSiteURL:        "abc",
				ReconciliationID:       reconciliationID,
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderSpecificData: map[string]interface{}{
					"provider data 1": "aaa",
					"provider data 2": 123,
					"provider data 3": map[string]interface{}{
						"provider data 4": "bbb",
					},
				},
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				COFTransactionIndicators: &zooz.COFTransactionIndicators{
					CardEntryMode:           "consent_transaction",
					COFConsentTransactionID: "",
				},
			},
			&zooz.ClientInfo{
				IPAddress: "95.173.136.70",
				UserAgent: "my-user-agent",
			},
		)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, authorizationCreated.ID)
			assert.NotEmpty(t, authorizationCreated.Created)
			assert.NotEmpty(t, authorizationCreated.PaymentMethod.FingerPrint)
			assert.NotEmpty(t, authorizationCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, authorizationCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, authorizationCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, authorizationCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Authorization{
				ID: authorizationCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Failed",
					Category:    "provider_error",
					SubCategory: "",
					Description: "Something went wrong on the provider's side.",
				},
				Amount:           5000,
				Created:          authorizationCreated.Created, // ignore
				ReconciliationID: reconciliationID,
				PaymentMethod: zooz.PaymentMethod{
					Href:               "", // TODO: empty because is not assigned to customer?
					Type:               "tokenized",
					TokenType:          string(token.TokenType),
					PassLuhnValidation: token.PassLuhnValidation,
					Token:              token.Token,
					Created:            json.Number(token.Created),
					Customer:           "",
					AdditionalDetails:  nil,
					BinNumber:          json.Number(token.BinNumber),
					Vendor:             token.Vendor,
					Issuer:             token.Issuer,
					CardType:           token.CardType,
					Level:              token.Level,
					CountryCode:        token.CountryCode,
					HolderName:         "name",
					ExpirationDate:     "",
					Last4Digits:        last4(cardNumber),
					ShippingAddress:    nil,
					BillingAddress:     nil,
					FingerPrint:        authorizationCreated.PaymentMethod.FingerPrint, // ignore
				},
				ThreeDSecureAttributes: nil,
				Installments:           nil,
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "102",
					Description:           "DECLINED.",
					RawResponse:           authorizationCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     "",
					TransactionID:         authorizationCreated.ProviderData.TransactionID, // ignore
					ExternalID:            authorizationCreated.ProviderData.ExternalID,    // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				ProviderSpecificData:       zooz.DecodedJSON{},                         // why empty?
				ProviderConfiguration:      authorizationCreated.ProviderConfiguration, // ignore
				OriginatingPurchaseCountry: "RUS",                                      // based on IP
				IPAddress:                  "95.173.136.70",
				Redirection:                nil,
				AdditionalDetails: zooz.AdditionalDetails{
					"auth detail 1": "value 1",
					"auth detail 2": "value 2",
				},
				DecisionEngineExecution: authorizationCreated.DecisionEngineExecution,
			}, authorizationCreated)
		})

		authorizationRetrieved, err := client.Authorization().Get(context.Background(), payment.ID, authorizationCreated.ID)
		require.NoError(t, err)
		must(t, func() {
			// empty on New, nil on Get
			assert.Nil(t, authorizationRetrieved.ProviderSpecificData)
			authorizationRetrieved.ProviderSpecificData = zooz.DecodedJSON{}

			// decision engine exectuion is not returned when retreving auth by id
			authorizationCreated.DecisionEngineExecution = zooz.DecisionEngineExecution{}

			assert.Equal(t, authorizationCreated, authorizationRetrieved)
		})
	})

	t.Run("idempotency", func(t *testing.T) {
		t.Parallel()

		idempotencyKey1 := randomString(32)
		idempotencyKey2 := randomString(32) // different key -> error
		token1, _ := PrepareToken(t, client)
		token2, _ := PrepareToken(t, client) // can't change token
		payment := PreparePayment(t, client, 5000, nil)

		authorization1, err := client.Authorization().New(
			context.Background(),
			idempotencyKey1,
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token1.Token,
					CreditCardCvv: token1.EncryptedCVV,
				},
			},
			nil,
		)
		require.NoError(t, err)

		authorization2, err := client.Authorization().New(
			context.Background(),
			idempotencyKey1,
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token1.Token,
					CreditCardCvv: token1.EncryptedCVV,
				},
			},
			nil,
		)
		require.NoError(t, err)
		require.Equal(t, authorization1, authorization2)

		authorization3, err := client.Authorization().New(
			context.Background(),
			idempotencyKey1,
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token2.Token, // can't change token
					CreditCardCvv: token2.EncryptedCVV,
				},
			},
			nil,
		)
		require.NoError(t, err)
		require.Equal(t, authorization1, authorization3)

		_, err = client.Authorization().New(
			context.Background(),
			idempotencyKey2, // different key
			payment.ID,
			&zooz.AuthorizationParams{
				PaymentMethod: zooz.PaymentMethodDetails{
					Type:          "tokenized",
					Token:         token1.Token,
					CreditCardCvv: token1.EncryptedCVV,
				},
			},
			nil,
		)
		requireZoozError(t, err, http.StatusConflict, zooz.APIError{
			Category:    "api_request_error",
			Description: "There was conflict with payment resource current state.",
			MoreInfo:    "Please check the current state of the payment.",
		})
	})

	t.Run("get - unknown authorization", func(t *testing.T) {
		t.Parallel()

		payment := PreparePayment(t, client, 123, nil)

		_, err := client.Authorization().Get(context.Background(), payment.ID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "authorizations resource does not exits",
		})
	})

	t.Run("get - unknown payment", func(t *testing.T) {
		t.Parallel()

		_, err := client.Authorization().Get(context.Background(), UnknownUUID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "Payment resource does not exists",
		})
	})
}

// PrepareAuthorization is a helper to create new authorization in zooz.
func PrepareAuthorization(t *testing.T, client *zooz.Client, payment *zooz.Payment, token *zooz.CreditCardToken) *zooz.Authorization {
	authorization, err := client.Authorization().New(
		context.Background(),
		randomString(32),
		payment.ID,
		&zooz.AuthorizationParams{
			PaymentMethod: zooz.PaymentMethodDetails{
				Type:          "tokenized",
				Token:         token.Token,
				CreditCardCvv: token.EncryptedCVV,
			},
			MerchantSiteURL:        "abc",
			ReconciliationID:       randomString(32),
			ThreeDSecureAttributes: nil,
			Installments:           nil,
			ProviderSpecificData: map[string]interface{}{
				"provider data 1": "aaa",
				"provider data 2": 123,
				"provider data 3": map[string]interface{}{
					"provider data 4": "bbb",
				},
			},
			AdditionalDetails: zooz.AdditionalDetails{
				"auth detail 1": "value 1",
				"auth detail 2": "value 2",
			},
			COFTransactionIndicators: &zooz.COFTransactionIndicators{
				CardEntryMode:           "consent_transaction",
				COFConsentTransactionID: "",
			},
		},
		&zooz.ClientInfo{
			IPAddress: "95.173.136.70",
			UserAgent: "my-user-agent",
		},
	)
	require.NoError(t, err)
	return authorization
}
