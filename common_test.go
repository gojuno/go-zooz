package zooz_test

import (
	"encoding/json"
	"testing"

	"github.com/gtforge/go-zooz"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestDecodeJSON__UnmarshalJSON(t *testing.T) {
	type testCase struct {
		name         string
		incomingData []byte
		expectedErr  error
		expectedRes  zooz.DecodedJSON
	}

	testCases := []testCase{
		{
			name: "positive",
			incomingData: []byte(`
				"{\"a\":\"{\\\"b\\\":1}\"}"
			`),
			expectedRes: zooz.DecodedJSON{
				"a": zooz.DecodedJSON{
					"b": float64(1),
				},
			},
		},
		{
			name:         "json is not correct",
			incomingData: []byte("asdasd"),
			expectedErr:  errors.New("invalid character 'a' looking for beginning of value"),
		},
		{
			name: "more complicated json",
			incomingData: []byte(`
			"{\"c\":\"{\\\"b\\\":\\\"{\\\\\\\"a\\\\\\\":1}\\\"}\"}"
		`),
			expectedRes: zooz.DecodedJSON{
				"c": zooz.DecodedJSON{
					"b": zooz.DecodedJSON{
						"a": float64(1),
					},
				},
			},
		},
		{
			name: "more more complicated json",
			incomingData: []byte(`
			"{\"c\":{\"b\":{\"c\":{\"e\":{\"a\":1}}}},\"d\":{\"a\":{\"1\":\"\\\"\\\\\\\"{\\\\\\\\\\\\\\\"b\\\\\\\\\\\\\\\":{\\\\\\\\\\\\\\\"c\\\\\\\\\\\\\\\":{\\\\\\\\\\\\\\\"e\\\\\\\\\\\\\\\":{\\\\\\\\\\\\\\\"a\\\\\\\\\\\\\\\":1}}}}\\\\\\\"\\\"\"}}}"
		`),
			expectedRes: zooz.DecodedJSON{
				"c": zooz.DecodedJSON{
					"b": zooz.DecodedJSON{
						"c": zooz.DecodedJSON{
							"e": zooz.DecodedJSON{
								"a": float64(1),
							},
						},
					},
				},
				"d": zooz.DecodedJSON{
					"a": zooz.DecodedJSON{
						"1": zooz.DecodedJSON{
							"b": zooz.DecodedJSON{
								"c": zooz.DecodedJSON{
									"e": zooz.DecodedJSON{
										"a": float64(1),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "do not need to unquote string",
			incomingData: []byte(`
			{"b":{"c":{"e":{"a":1}}}}
			`),
			expectedRes: zooz.DecodedJSON{
				"b": zooz.DecodedJSON{
					"c": zooz.DecodedJSON{
						"e": zooz.DecodedJSON{
							"a": float64(1),
						},
					},
				},
			},
		},
		{
			name: "when json contains quoted array",
			incomingData: []byte(`
			{"c":{"b":[2,3,4,5]}}
			`),
			expectedRes: zooz.DecodedJSON{
				"c": zooz.DecodedJSON{
					"b": []interface{}{
						float64(2),
						float64(3),
						float64(4),
						float64(5),
					},
				},
			},
		},
		{
			name: "more different types in JSON",
			incomingData: []byte(`
			{"a":1,"b":"2021-04-14 18:32:30 +0300","c":12.12,"d":["foo"],"e":{"k":1},"foo":"asd"}
			`),
			expectedRes: zooz.DecodedJSON{
				"a": float64(1),
				"b": "2021-04-14 18:32:30 +0300",
				"c": float64(12.12),
				"d": []interface{}{"foo"},
				"e": zooz.DecodedJSON{
					"k": float64(1),
				},
				"foo": "asd",
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			res := zooz.DecodedJSON{}
			err := json.Unmarshal(tC.incomingData, &res)
			if tC.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tC.expectedErr.Error())
				return
			}
			assert.Equal(t, tC.expectedRes, res)
		})
	}
}

const authorizationRawResponse = `
{
	"id": "a3cf4923-582b-4e00-ba49-705fea6753ae",
	"created": "1618824537698",
	"reconciliation_id": "1512196148",
	"payment_method": {
		"type": "tokenized",
		"token": "c8969685-d93a-4007-8ce8-2c05548b36a6",
		"token_type": "credit_card",
		"additional_details": {
			"app_platform": "IPHONE",
			"app_version": "9.62.1",
			"user_id": "1234567"
		},
		"holder_name": "Debbie Roberts",
		"expiration_date": "04/2025",
		"last_4_digits": "9441",
		"pass_luhn_validation": true,
		"fingerprint": "b57eda35-758d-425b-9345-eb9c8449dc3c",
		"bin_number": "98123",
		"vendor": "VISA",
		"issuer": "SBERBANK",
		"card_type": "DEBIT",
		"level": "CLASSIC",
		"country_code": "GBR",
		"created": "1594130055674"
	},
	"ip_address": "0cf4:bc67:8d4a:6e66:b878:b843:ed50:a123",
	"originating_purchase_country": "GBR",
	"additional_details": {
		"dunning_count": "0",
		"payment_created_by": "authorization"
	},
	"result": {
		"status": "Succeed"
	},
	"provider_specific_data": {
		"Bank": "Sberbank"
	},
	"provider_data": {
		"provider_name": "SafeCharge",
		"response_code": "APPROVED",
		"raw_response": "{\"provider_response\":\"{\\\"threeDFlow\\\":\\\"0\\\",\\\"orderId\\\":12345667,\\\"transactionStatus\\\":\\\"APPROVED\\\",\\\"gwErrorCode\\\":0,\\\"gwExtendedErrorCode\\\":0,\\\"userPaymentOptionId\\\":\\\"123212\\\",\\\"externalTransactionId\\\":\\\"\\\",\\\"transactionId\\\":\\\"1131234561698168588\\\",\\\"authCode\\\":\\\"004037\\\",\\\"userTokenId\\\":\\\"3098afd5-65eb-470a-9644-7bf48dc7795e\\\",\\\"CVV2Reply\\\":\\\"\\\",\\\"AVSCode\\\":\\\"N\\\",\\\"customData\\\":\\\"\\\",\\\"transactionType\\\":\\\"Auth\\\",\\\"fraudDetails\\\":{\\\"finalDecision\\\":\\\"Accept\\\",\\\"score\\\":\\\"0\\\"},\\\"acquirerId\\\":\\\"103\\\",\\\"bin\\\":\\\"123456\\\",\\\"last4Digits\\\":\\\"1231\\\",\\\"ccCardNumber\\\":\\\"1****1234\\\",\\\"ccExpMonth\\\":\\\"04\\\",\\\"ccExpYear\\\":\\\"25\\\",\\\"cardType\\\":\\\"Debit\\\",\\\"cardBrand\\\":\\\"VISA\\\",\\\"sessionToken\\\":\\\"251151f3-444e-42f9-9d7b-9d4da2c36d67\\\",\\\"clientUniqueId\\\":\\\"12345678\\\",\\\"internalRequestId\\\":12345677,\\\"status\\\":\\\"SUCCESS\\\",\\\"errCode\\\":0,\\\"reason\\\":\\\"\\\",\\\"merchantId\\\":\\\"5655545467654\\\",\\\"merchantSiteId\\\":\\\"1212121\\\",\\\"version\\\":\\\"1.0\\\",\\\"clientRequestId\\\":\\\"ad0107d1-14f8-41e9-b45c-dd5ec87d91b6\\\"}\"}",
		"authorization_code": "004037",
		"transaction_id": "1130000001698168588",
		"external_id": "19001609518"
	},
	"amount": 100,
	"provider_configuration": {
		"id": "b7ff80dc-b56c-403a-8036-295e1aebd5b7",
		"name": "SafeCharge",
		"created": "1534318898741",
		"modified": "1534318898741",
		"account_id": "2058ac1f-ca4b-497f-8181-341e0eea5392",
		"provider_id": "d54f3610-3722-4d76-a785-bfbcfdd173dd",
		"type": "cc_processor",
		"href": "https://api.paymentsos.com/accounts/2058ac1f-ca4b-497f-8181-341e0eea5392/provider-configurations/b7ff80dc-b56c-403a-8036-295e1aebd5b7"
	},
	"decision_engine_execution": {
		"id": "8bc51c25-8f9a-4ff2-a065-e3348c35f495",
		"flow_id": "8bc51c25-8f9a-4ff2-a065-e3348c35f495",
		"created": "2021-04-19T09:28:56.748Z",
		"status": "Completed",
		"policy_results": [
			{
				"type": "SelectionPolicy",
				"name": "VT",
				"execution_time": "2021-04-19T09:28:56.760267Z",
				"result": "Skip",
				"selection_skipped_reason": "NoMatchingConditions"
			},
			{
				"type": "TargetPolicy",
				"name": "TargetPolicy",
				"execution_time": "2021-04-19T09:28:57.748587Z",
				"result": "Hit",
				"provider_name": "SafeCharge",
				"provider_configuration": "https://api.paymentsos.com/accounts/4e0fdf9f-4080-41d3-ac10-84defd5aeca2/providers-configurations/4e0fdf9f-4080-41d3-ac10-84defd5aeca2",
				"transaction": "https://api.paymentsos.com/payments/4e0fdf9f-4080-41d3-ac10-84defd5aeca2/authorizations/4e0fdf9f-4080-41d3-ac10-84defd5aeca2"
			}
		]
	}
	}`

func TestUnmarshalAuthorization(t *testing.T) {
	authorization := zooz.Authorization{}
	err := json.Unmarshal([]byte(authorizationRawResponse), &authorization)
	assert.NoError(t, err)
	assert.Equal(t, zooz.Authorization{
		ID: "a3cf4923-582b-4e00-ba49-705fea6753ae",
		Result: zooz.Result{
			Status:      "Succeed",
			Category:    "",
			SubCategory: "",
			Description: "",
		},
		Amount:           100,
		Created:          "1618824537698",
		ReconciliationID: "1512196148",
		PaymentMethod: zooz.PaymentMethod{
			Href:               "",
			Type:               "tokenized",
			TokenType:          "credit_card",
			PassLuhnValidation: true,
			Token:              "c8969685-d93a-4007-8ce8-2c05548b36a6",
			Created:            "1594130055674",
			Customer:           "",
			AdditionalDetails: zooz.AdditionalDetails{
				"app_platform": "IPHONE",
				"app_version":  "9.62.1",
				"user_id":      "1234567",
			},
			BinNumber:      "98123",
			Vendor:         "VISA",
			Issuer:         "SBERBANK",
			CardType:       "DEBIT",
			Level:          "CLASSIC",
			CountryCode:    "GBR",
			HolderName:     "Debbie Roberts",
			ExpirationDate: "04/2025",
			Last4Digits:    "9441",
			FingerPrint:    "b57eda35-758d-425b-9345-eb9c8449dc3c",
		},
		ProviderData: zooz.ProviderData{
			ProviderName: "SafeCharge",
			ResponseCode: "APPROVED",
			Description:  "",
			RawResponse: zooz.DecodedJSON{
				"provider_response": zooz.DecodedJSON{
					"AVSCode":               "N",
					"CVV2Reply":             "",
					"acquirerId":            "103",
					"authCode":              "004037",
					"bin":                   "123456",
					"cardBrand":             "VISA",
					"cardType":              "Debit",
					"ccCardNumber":          "1****1234",
					"ccExpMonth":            "04",
					"ccExpYear":             "25",
					"clientRequestId":       "ad0107d1-14f8-41e9-b45c-dd5ec87d91b6",
					"clientUniqueId":        "12345678",
					"customData":            "",
					"errCode":               float64(0),
					"externalTransactionId": "",
					"fraudDetails": zooz.DecodedJSON{
						"finalDecision": "Accept",
						"score":         "0",
					},
					"gwErrorCode":         float64(0),
					"gwExtendedErrorCode": float64(0),
					"internalRequestId":   1.2345677e+07,
					"last4Digits":         "1231",
					"merchantId":          "5655545467654",
					"merchantSiteId":      "1212121",
					"orderId":             1.2345667e+07,
					"reason":              "",
					"sessionToken":        "251151f3-444e-42f9-9d7b-9d4da2c36d67",
					"status":              "SUCCESS",
					"threeDFlow":          "0",
					"transactionId":       "1131234561698168588",
					"transactionStatus":   "APPROVED",
					"transactionType":     "Auth",
					"userPaymentOptionId": "123212",
					"userTokenId":         "3098afd5-65eb-470a-9644-7bf48dc7795e",
					"version":             "1.0",
				},
			},
			AvsCode:           "",
			AuthorizationCode: "004037",
			TransactionID:     "1130000001698168588",
			ExternalID:        "19001609518",
		},
		ProviderSpecificData: zooz.DecodedJSON{
			"Bank": "Sberbank",
		},
		ProviderConfiguration: zooz.ProviderConfiguration{
			ID:          "b7ff80dc-b56c-403a-8036-295e1aebd5b7",
			Name:        "SafeCharge",
			Description: "",
			Created:     "1534318898741",
			Modified:    "1534318898741",
			ProviderID:  "d54f3610-3722-4d76-a785-bfbcfdd173dd",
			Type:        "cc_processor",
			AccountID:   "2058ac1f-ca4b-497f-8181-341e0eea5392",
			Href:        "https://api.paymentsos.com/accounts/2058ac1f-ca4b-497f-8181-341e0eea5392/provider-configurations/b7ff80dc-b56c-403a-8036-295e1aebd5b7",
		},
		OriginatingPurchaseCountry: "GBR",
		IPAddress:                  "0cf4:bc67:8d4a:6e66:b878:b843:ed50:a123",
		AdditionalDetails: zooz.AdditionalDetails{
			"dunning_count":      "0",
			"payment_created_by": "authorization",
		},
		DecisionEngineExecution: zooz.DecisionEngineExecution{
			ID:      "8bc51c25-8f9a-4ff2-a065-e3348c35f495",
			Created: "2021-04-19T09:28:56.748Z",
			FlowID:  "8bc51c25-8f9a-4ff2-a065-e3348c35f495",
			Status:  "Completed",
			PolicyResults: []zooz.PolicyResult{
				{
					Type:                  "SelectionPolicy",
					ProviderName:          "",
					ProviderConfiguration: "",
					Name:                  "VT",
					ExecutionTime:         "2021-04-19T09:28:56.760267Z",
					Transaction:           "",
					Result:                "Skip",
				},
				{
					Type:                  "TargetPolicy",
					ProviderName:          "SafeCharge",
					ProviderConfiguration: "https://api.paymentsos.com/accounts/4e0fdf9f-4080-41d3-ac10-84defd5aeca2/providers-configurations/4e0fdf9f-4080-41d3-ac10-84defd5aeca2",
					Name:                  "TargetPolicy",
					ExecutionTime:         "2021-04-19T09:28:57.748587Z",
					Transaction:           "https://api.paymentsos.com/payments/4e0fdf9f-4080-41d3-ac10-84defd5aeca2/authorizations/4e0fdf9f-4080-41d3-ac10-84defd5aeca2",
					Result:                "Hit",
				},
			},
		}}, authorization)
}
