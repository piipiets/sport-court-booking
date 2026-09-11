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

// @Summary      Record a payment
// @Description  Record a payment for the authenticated user's booking
// @Tags         Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body request.CreatePaymentRequest true "Payment payload"
// @Success      200 {object} common.APIResponse
// @Failure      400 {object} common.APIResponse
// @Failure      401 {object} common.APIResponse
// @Failure      403 {object} common.APIResponse
// @Failure      404 {object} common.APIResponse
// @Failure      409 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/payments [post]
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

// @Summary      Get payment by booking ID
// @Description  Retrieve the payment for a specific booking (owner or admin)
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        booking_id path int true "Booking ID"
// @Success      200 {object} common.APIResponse{data=response.PaymentResponse}
// @Failure      400 {object} common.APIResponse
// @Failure      401 {object} common.APIResponse
// @Failure      403 {object} common.APIResponse
// @Failure      404 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/payments/{booking_id} [get]
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

// @Summary      Get all payments for the logged-in user
// @Description  Retrieve all payments for the authenticated user (admins see all)
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} common.APIResponse{data=[]response.PaymentResponse}
// @Failure      401 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/payments [get]
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
