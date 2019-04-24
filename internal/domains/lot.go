package domains

import "time"

type LotGeneral struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	BuyPrice    *float64   `json:"buy_price,omitempty"`
	MinPrice    float64    `json:"min_price"`
	PriceStep   float64    `json:"price_step"`
	Status      string     `json:"status"`
	EndAt       time.Time  `json:"end_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
type Lot struct {
	LotGeneral
	CreatorID int64  `json:"creator"`
	BuyerID   *int64 `json:"buyer,omitempty"`
}

func (lg *Lot) IsDeleted() bool {
	return lg.DeletedAt != nil

}
