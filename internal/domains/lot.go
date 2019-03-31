package domains

import (
	"time"
)

type Lot struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description, omitempty"`
	MinPrice    float64   `json:"min_price"`
	PriceStep   float64   `json:"price_step"`
	Status      string    `json:"status"`
	EndAt       time.Time `json:"end_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatorID   int       `json:"-"`
	BuyerID     int       `json:"-"`
}
