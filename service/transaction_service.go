package service

import (
	"context"
	"log"

	"github.com/Shabrinashsf/go-xendit-payment-webhook/constants"
	"github.com/Shabrinashsf/go-xendit-payment-webhook/dto"
	"github.com/Shabrinashsf/go-xendit-payment-webhook/entity"
	"github.com/Shabrinashsf/go-xendit-payment-webhook/repository"
	"github.com/Shabrinashsf/go-xendit-payment-webhook/utils/payment"
	"github.com/google/uuid"
)

type (
	TransactionService interface {
		CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (dto.CreateTransactionResponse, error)
		XenditWebhook(ctx context.Context, req dto.XenditWebhook) error
	}

	transactionService struct {
		transactionRepo repository.TransactionRepository
	}
)

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (dto.CreateTransactionResponse, error) {
	product, err := s.transactionRepo.GetProductByID(ctx, nil, uuid.MustParse(req.ProductID))
	if err != nil {
		return dto.CreateTransactionResponse{}, err
	}

	item := dto.XenditItem{
		Name:     product.Name,
		Price:    int(product.Price),
		Quantity: 1,
	}

	customer := dto.XenditCustomer{
		GivenNames:   req.Name,
		Email:        req.Email,
		MobileNumber: req.MobileNumber,
	}

	transactionId := uuid.New()

	invoice := dto.XenditInvoice{
		ExternalID:         transactionId.String(),
		Amount:             product.Price,
		Description:        product.Name,
		InvoiceDuration:    constants.INVOICE_DURATION,
		Customer:           customer,
		SuccessRedirectURL: constants.SUCCESS_REDIRECT_URL,
		FailureRedirectURL: constants.FAILED_REDIRECT_URL,
		Currency:           "IDR",
		Items:              []dto.XenditItem{item},
	}

	res, err := payment.SendXenditInvoice(invoice)
	if err != nil {
		return dto.CreateTransactionResponse{}, dto.ErrFailedCreateInvoice
	}

	log.Println(res)

	invoiceURL, ok := res["invoice_url"].(string)
	if !ok {
		return dto.CreateTransactionResponse{}, dto.ErrFailedCreateInvoice
	}

	transaction := entity.Transaction{
		ID:         transactionId,
		ProductID:  product.ID,
		AmountPaid: 0,
		Status:     "PENDING",
		InvoiceUrl: invoiceURL,
	}

	trans, err := s.transactionRepo.CreateTransaction(ctx, nil, transaction)
	if err != nil {
		return dto.CreateTransactionResponse{}, err
	}

	return dto.CreateTransactionResponse{
		InvoiceURL: trans.InvoiceUrl,
	}, nil
}

func (s *transactionService) XenditWebhook(ctx context.Context, req dto.XenditWebhook) error {
	id, err := uuid.Parse(req.ExternalID)
	if err != nil {
		return dto.ErrParseUUID
	}

	transaction, err := s.transactionRepo.GetTransactionByID(ctx, nil, id)
	if err != nil {
		return err
	}

	switch req.Status {
	case "PAID", "SETTLED":
		transaction.Status = req.Status
		transaction.AmountPaid = req.Amount
		if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
			return err
		}
	case "CANCELLED":
		transaction.Status = "CANCELLED"
		if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
			return err
		}
		if err := s.transactionRepo.DeleteTransaction(ctx, nil, transaction.ID); err != nil {
			return err
		}
	case "EXPIRED":
		transaction.Status = "EXPIRED"
		if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
			return err
		}
		if err := s.transactionRepo.DeleteTransaction(ctx, nil, transaction.ID); err != nil {
			return err
		}
	default:
		return dto.ErrStatusUnknownPayment
	}

	return nil
}
