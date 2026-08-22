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

type CourtHandler struct {
	courtService service.CourtService
}

func NewCourtHandler(courtService service.CourtService) *CourtHandler {
	return &CourtHandler{courtService: courtService}
}

func (h *CourtHandler) GetAll(c *gin.Context) {
	courts, err := h.courtService.GetAll()
	if err != nil {
		fmt.Println("Error : ", err)
		common.GenerateErrorResponse(c, "failed to fetch courts", http.StatusBadRequest)
		return
	}

	common.GenerateSuccessResponseWithData(c, "fetched courts successfully", courts)
}

func (h *CourtHandler) GetByID(c *gin.Context) {
	id, err := parseCourtID(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	court, err := h.courtService.GetByID(id)
	if errors.Is(err, constant.ErrCourtNotFound) {
		common.GenerateErrorResponse(c, "court not found", http.StatusNotFound)
		return
	}
	if err != nil {
		common.GenerateErrorResponse(c, "failed to fetch court", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponseWithData(c, "fetched court successfully", court)
}

func (h *CourtHandler) Create(c *gin.Context) {
	var req request.CourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.courtService.Create(req)
	if err != nil {
		common.GenerateErrorResponse(c, "failed to create court", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponse(c, "court created successfully")
}

func (h *CourtHandler) Update(c *gin.Context) {
	id, err := parseCourtID(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadGateway)
		return
	}

	var req request.UpdateCourtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.courtService.Update(id, req)
	if errors.Is(err, constant.ErrCourtNotFound) {
		common.GenerateErrorResponse(c, "court not found", http.StatusNotFound)
		return
	}

	if err != nil {
		common.GenerateErrorResponse(c, "failed to update court", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponse(c, "court updated successfully")
}

func (h *CourtHandler) Delete(c *gin.Context) {
	id, err := parseCourtID(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.courtService.Delete(id)
	if errors.Is(err, constant.ErrCourtNotFound) {
		common.GenerateErrorResponse(c, "court not found", http.StatusNotFound)
		return
	}

	if err != nil {
		common.GenerateErrorResponse(c, "failed to delete court", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponse(c, "court deleted")
}

func parseCourtID(c *gin.Context) (int64, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid court id")
	}
	return id, nil
}
