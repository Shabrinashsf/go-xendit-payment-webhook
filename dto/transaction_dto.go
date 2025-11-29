package dto

import "errors"

const (
	MESSAGE_SUCCESS_CREATE_TRANSACTION  = "success create transaction"
	MESSAGE_SUCCESS_GET_CALLBACK_XENDIT = "success get callback from xendit"

	MESSAGE_FAILED_GET_DATA_FROM_BODY  = "failed to get data from body"
	MESSAGE_FAILED_CREATE_TRANSACTION  = "failed to create transaction"
	MESSAGE_FAILED_GET_CALLBACK_XENDIT = "failed to get callback from xendit"
)

var (
	ErrFailedCreateInvoice  = errors.New("failed to create invoice")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrParseUUID            = errors.New("failed to parse uuid")
	ErrStatusUnknownPayment = errors.New("status unknown payment")
)

type (
	XenditCustomer struct {
		GivenNames   string `json:"given_names"`
		Email        string `json:"email"`
		MobileNumber string `json:"mobile_number"`
	}

	XenditItem struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Price    int    `json:"price"`
	}

	XenditInvoice struct {
		ExternalID         string         `json:"external_id"`
		Amount             int            `json:"amount"`
		Description        string         `json:"description"`
		InvoiceDuration    int            `json:"invoice_duration"`
		Customer           XenditCustomer `json:"customer"`
		SuccessRedirectURL string         `json:"success_redirect_url"`
		FailureRedirectURL string         `json:"failure_redirect_url"`
		Currency           string         `json:"currency"`
		Items              []XenditItem   `json:"items"`
	}

	CreateTransactionRequest struct {
		Name         string `json:"name" binding:"required"`
		Email        string `json:"email" binding:"required,email"`
		MobileNumber string `json:"mobile_number" binding:"required"`
		ProductID    string `json:"product_id" binding:"required,uuid"`
	}

	CreateTransactionResponse struct {
		InvoiceURL string `json:"invoice_url"`
	}

	XenditWebhook struct {
		ID                       string `json:"id"`
		ExternalID               string `json:"external_id"`
		UserID                   string `json:"user_id"`
		IsHigh                   bool   `json:"is_high"`
		PaymentMethod            string `json:"payment_method"`
		Status                   string `json:"status"`
		MerchantName             string `json:"merchant_name"`
		Amount                   int    `json:"amount"`
		PaidAmount               int    `json:"paid_amount"`
		BankCode                 string `json:"bank_code"`
		PaidAt                   string `json:"paid_at"`
		PayerEmail               string `json:"payer_email"`
		Description              string `json:"description"`
		AdjustmentReceivedAmount int    `json:"adjustment_received_amount"`
		FeesPaidAmount           int    `json:"fees_paid_amount"`
		Updated                  string `json:"updated"`
		Created                  string `json:"created"`
		Currency                 string `json:"currency"`
		PaymentChannel           string `json:"payment_channel"`
		PaymentDestination       string `json:"payment_destination"`
	}
)
