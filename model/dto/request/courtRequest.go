package request

type CourtRequest struct {
	Name     string  `json:"name" binding:"required,min=3,max=150"`
	Type     string  `json:"type" binding:"required,oneof=futsal badminton"`
	Price    float64 `json:"price_per_hour" binding:"required,gt=0"`
	Location string  `json:"location" binding:"required,min=3,max=255"`
}

// for partially update
type UpdateCourtRequest struct {
	Name     *string  `json:"name" binding:"omitempty,min=3,max=150"`
	Type     *string  `json:"type" binding:"omitempty,oneof=futsal badminton"`
	Price    *float64 `json:"price" binding:"omitempty,gt=0"`
	Location *string  `json:"location" binding:"omitempty,min=3,max=255"`
}
