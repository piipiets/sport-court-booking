package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/helpers/common"
	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/middlewares"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/service"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// POST /api/bookings
func (h *BookingHandler) Create(c *gin.Context) {
	var req request.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusUnauthorized)
		return
	}

	err = h.bookingService.Create(userID, req)
	if errors.Is(err, constant.ErrBookingConflict) {
		common.GenerateErrorResponse(c, err.Error(), http.StatusConflict)
		return
	}
	if errors.Is(err, constant.ErrCourtNotFound) {
		common.GenerateErrorResponse(c, "court not found", http.StatusNotFound)
		return
	}
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponse(c, "booking created successfully")
}

// GET /api/bookings — booking milik user yang login
func (h *BookingHandler) GetAll(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	bookings, err := h.bookingService.GetByUserID(userID)
	if err != nil {
		common.GenerateErrorResponse(c, "failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponseWithData(c, "success", bookings)
}

// GET /api/bookings/:id
func (h *BookingHandler) GetByID(c *gin.Context) {
	id, err := parseBookingID(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	isAdmin := isAdminFromContext(c)

	booking, err := h.bookingService.GetByID(id, userID, isAdmin)
	if errors.Is(err, constant.ErrBookingNotFound) {
		common.GenerateErrorResponse(c, "booking not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, constant.ErrForbidden) {
		common.GenerateErrorResponse(c, err.Error(), http.StatusForbidden)
		return
	}
	if err != nil {
		common.GenerateErrorResponse(c, "failed to fetch booking", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponseWithData(c, "success", booking)
}

// PUT /api/bookings/:id/status — admin only
func (h *BookingHandler) UpdateStatus(c *gin.Context) {
	id, err := parseBookingID(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	role := isAdminFromContext(c)
	if !role {
		common.GenerateErrorResponse(c, "No Access", http.StatusForbidden)
		return
	}

	var req request.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.bookingService.UpdateStatus(id, req)
	if errors.Is(err, constant.ErrBookingNotFound) {
		common.GenerateErrorResponse(c, "booking not found", http.StatusNotFound)
		return
	}
	if err != nil {
		common.GenerateErrorResponse(c, "failed to update booking status", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponse(c, "booking status updated")
}

func parseBookingID(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid booking id")
	}
	return id, nil
}

func getAuthClaims(c *gin.Context) (*middlewares.Claims, error) {
	authRaw, exists := c.Get("auth")
	if !exists {
		return nil, errors.New("unauthorized")
	}

	claims, ok := authRaw.(*middlewares.Claims)
	if !ok {
		return nil, errors.New("invalid auth context")
	}
	return claims, nil
}

func getUserIDFromContext(c *gin.Context) (int64, error) {
	claims, err := getAuthClaims(c)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func isAdminFromContext(c *gin.Context) bool {
	claims, err := getAuthClaims(c)
	if err != nil {
		return false
	}
	return claims.Role == "admin"
}
