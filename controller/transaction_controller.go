package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/Shabrinashsf/go-xendit-payment-webhook/dto"
	"github.com/Shabrinashsf/go-xendit-payment-webhook/service"
	"github.com/Shabrinashsf/go-xendit-payment-webhook/utils/response"
	"github.com/gin-gonic/gin"
)

type (
	TransactionController interface {
		CreateTransaction(ctx *gin.Context)
		XenditWebhook(ctx *gin.Context)
	}

	transactionController struct {
		transactionService service.TransactionService
	}
)

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

func (c *transactionController) CreateTransaction(ctx *gin.Context) {
	// Dont forget to implement your business logic
	var req dto.CreateTransactionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.transactionService.CreateTransaction(ctx.Request.Context(), req)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TRANSACTION, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_TRANSACTION, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *transactionController) XenditWebhook(ctx *gin.Context) {
	svcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var req dto.XenditWebhook
	if err := ctx.ShouldBind(&req); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.transactionService.XenditWebhook(svcCtx, req); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CALLBACK_XENDIT, err.Error(), nil)
		switch err {
		case dto.ErrTransactionNotFound:
			ctx.JSON(http.StatusNotFound, res)
		case dto.ErrParseUUID, dto.ErrStatusUnknownPayment:
			ctx.JSON(http.StatusBadRequest, res)
		default:
			ctx.JSON(http.StatusInternalServerError, res)
		}
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_CALLBACK_XENDIT, nil)
	ctx.JSON(http.StatusOK, res)
}
