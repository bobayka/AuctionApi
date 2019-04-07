package domains

import (
	"gitlab.com/bobayka/courseproject/internal/responce"
	"time"
)

type Lot struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description *string            `json:"description,omitempty"`
	BuyPrice    *float64           `json:"buy_price"`
	MinPrice    float64            `json:"min_price"`
	PriceStep   float64            `json:"price_step"`
	Status      string             `json:"status"`
	EndAt       time.Time          `json:"end_at"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CreatorID   responce.ShortUSer `json:"creator"`
	BuyerID     responce.ShortUSer `json:"buyer"`
}
