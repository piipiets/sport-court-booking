package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/helpers/common"
	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/service"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// POST /api/payments
func (h *PaymentHandler) Create(c *gin.Context) {
	var req request.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.paymentService.Create(userID, req)
	switch {
	case errors.Is(err, constant.ErrBookingNotFound):
		common.GenerateErrorResponse(c, "booking not found", http.StatusNotFound)
		return
	case errors.Is(err, constant.ErrForbidden):
		common.GenerateErrorResponse(c, err.Error(), http.StatusForbidden)
		return
	case errors.Is(err, constant.ErrPaymentAmountMismatch):
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, constant.ErrPaymentAlreadyExists):
		common.GenerateErrorResponse(c, err.Error(), http.StatusConflict)
		return
	case err != nil:
		common.GenerateErrorResponse(c, "failed to create payment", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponse(c, "payment recorded successfully")
}

// GET /api/payments/:booking_id
func (h *PaymentHandler) GetByBookingID(c *gin.Context) {
	bookingID, err := strconv.ParseInt(c.Param("booking_id"), 10, 64)
	if err != nil {
		common.GenerateErrorResponse(c, "invalid booking id", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusUnauthorized)
		return
	}
	isAdmin := isAdminFromContext(c)

	payment, err := h.paymentService.GetByBookingID(bookingID, userID, isAdmin)
	switch {
	case errors.Is(err, constant.ErrBookingNotFound):
		common.GenerateErrorResponse(c, "booking not found", http.StatusNotFound)
		return
	case errors.Is(err, constant.ErrForbidden):
		common.GenerateErrorResponse(c, err.Error(), http.StatusForbidden)
		return
	case errors.Is(err, constant.ErrPaymentNotFound):
		common.GenerateErrorResponse(c, "payment not found for this booking", http.StatusNotFound)
		return
	case err != nil:
		common.GenerateErrorResponse(c, "failed to fetch payment", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponseWithData(c, "success", payment)
}

// GET /api/payments/:booking_id
func (h *PaymentHandler) GetAllPaymentByUserId(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusUnauthorized)
		return
	}

	isAdmin := isAdminFromContext(c)

	payment, err := h.paymentService.GetAllPaymentsByUserID(isAdmin, userID)

	if err != nil {
		fmt.Println(err)
		common.GenerateErrorResponse(c, "failed to fetch payment", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponseWithData(c, "success", payment)
}
