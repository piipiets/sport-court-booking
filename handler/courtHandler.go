package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/piipiets/sport-court-booking/helpers/common"
	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/service"
)

const maxImageSize = 5 * 1024 * 1024

type CourtHandler struct {
	courtService service.CourtService
}

func NewCourtHandler(courtService service.CourtService) *CourtHandler {
	return &CourtHandler{courtService: courtService}
}

// @Summary      Get all courts
// @Description  Fetch all sport courts
// @Tags         Courts
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} common.APIResponse{data=[]response.CourtResponse}
// @Failure      400 {object} common.APIResponse
// @Router       /api/courts [get]
func (h *CourtHandler) GetAll(c *gin.Context) {
	courts, err := h.courtService.GetAll()
	if err != nil {
		fmt.Println("Error : ", err)
		common.GenerateErrorResponse(c, "failed to fetch courts", http.StatusBadRequest)
		return
	}

	common.GenerateSuccessResponseWithData(c, "fetched courts successfully", courts)
}

// @Summary      Get a court by ID
// @Description  Fetch a single court by its ID
// @Tags         Courts
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Court ID"
// @Success      200 {object} common.APIResponse{data=response.CourtResponse}
// @Failure      400 {object} common.APIResponse
// @Failure      404 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/courts/{id} [get]
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

// @Summary      Get court image
// @Description  Stream the court image bytes
// @Tags         Courts
// @Security     BearerAuth
// @Produce      image/jpeg,image/png,image/webp
// @Param        id path int true "Court ID"
// @Success      200 {file} binary
// @Failure      400 {object} common.APIResponse
// @Failure      404 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/courts/{id}/image [get]
func (h *CourtHandler) GetImage(c *gin.Context) {
	id, err := parseCourtID(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	data, contentType, err := h.courtService.GetImage(id)
	if errors.Is(err, constant.ErrCourtNotFound) {
		common.GenerateErrorResponse(c, "court not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, constant.ErrCourtImageNotFound) {
		common.GenerateErrorResponse(c, "court image not found", http.StatusNotFound)
		return
	}
	if err != nil {
		common.GenerateErrorResponse(c, "failed to fetch court image", http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, contentType, data)
}

// @Summary      Create a court
// @Description  Create a new court (admin only)
// @Tags         Courts
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        name formData string true "Court name"
// @Param        type formData string true "Court type (futsal or badminton)"
// @Param        price_per_hour formData number true "Price per hour"
// @Param        location formData string true "Court location"
// @Param        image formData file false "Court image (jpg/jpeg/png/webp, max 5MB)"
// @Success      200 {object} common.APIResponse
// @Failure      400 {object} common.APIResponse
// @Failure      403 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/courts [post]
func (h *CourtHandler) Create(c *gin.Context) {
	if !isAdminFromContext(c) {
		common.GenerateErrorResponse(c, "No Access", http.StatusForbidden)
		return
	}

	var req request.CourtRequest
	if err := c.ShouldBind(&req); err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	image, err := parseUploadImage(c)
	if err != nil {
		fmt.Printf("Error parsing image: %v\n", err)
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.courtService.Create(req, image)
	if err != nil {
		fmt.Printf("Error creating court: %v\n", err)
		common.GenerateErrorResponse(c, "failed to create court", http.StatusInternalServerError)
		return
	}

	common.GenerateSuccessResponse(c, "court created successfully")
}

// @Summary      Update a court
// @Description  Update a court (admin only)
// @Tags         Courts
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path int true "Court ID"
// @Param        name formData string false "Court name"
// @Param        type formData string false "Court type (futsal or badminton)"
// @Param        price_per_hour formData number false "Price per hour"
// @Param        location formData string false "Court location"
// @Param        image formData file false "New court image"
// @Success      200 {object} common.APIResponse
// @Failure      400 {object} common.APIResponse
// @Failure      403 {object} common.APIResponse
// @Failure      404 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/courts/{id} [put]
func (h *CourtHandler) Update(c *gin.Context) {
	if !isAdminFromContext(c) {
		common.GenerateErrorResponse(c, "No Access", http.StatusForbidden)
		return
	}

	id, err := parseCourtID(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadGateway)
		return
	}

	var req request.UpdateCourtRequest
	if err := c.ShouldBind(&req); err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	image, err := parseUploadImage(c)
	if err != nil {
		common.GenerateErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.courtService.Update(id, req, image)
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

func parseUploadImage(c *gin.Context) (*service.UploadImage, error) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, nil
		}
		return nil, err
	}

	if err := validateImage(fileHeader); err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return &service.UploadImage{
		File:     file,
		FileName: fileHeader.Filename,
		Size:     fileHeader.Size,
	}, nil
}

func validateImage(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > maxImageSize {
		return errors.New("image must be 5MB or smaller")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return nil
	default:
		return errors.New("unsupported image type, only jpg, jpeg, png, webp are allowed")
	}
}

// @Summary      Delete a court
// @Description  Delete a court by ID (admin only)
// @Tags         Courts
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Court ID"
// @Success      200 {object} common.APIResponse
// @Failure      400 {object} common.APIResponse
// @Failure      403 {object} common.APIResponse
// @Failure      404 {object} common.APIResponse
// @Failure      500 {object} common.APIResponse
// @Router       /api/courts/{id} [delete]
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
