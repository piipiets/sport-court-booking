package request

type CourtRequest struct {
	Name     string  `form:"name" json:"name" binding:"required,min=3,max=150"`
	Type     string  `form:"type" json:"type" binding:"required,oneof=futsal badminton"`
	Price    float64 `form:"price_per_hour" json:"price_per_hour" binding:"required,gt=0"`
	Location string  `form:"location" json:"location" binding:"required,min=3,max=255"`
}

// for partially update
type UpdateCourtRequest struct {
	Name     *string  `form:"name" json:"name" binding:"omitempty,min=3,max=150"`
	Type     *string  `form:"type" json:"type" binding:"omitempty,oneof=futsal badminton"`
	Price    *float64 `form:"price_per_hour" json:"price_per_hour" binding:"omitempty,gt=0"`
	Location *string  `form:"location" json:"location" binding:"omitempty,min=3,max=255"`
}
