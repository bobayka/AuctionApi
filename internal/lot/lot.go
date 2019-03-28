package lot

import "gitlab.com/bobayka/courseproject/pkg/customTime"

type Lot struct {
	ID          int                    `json:"id"`
	CreatorID   int                    `json:"creator_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description, omitempty"`
	MinPrice    float64                `json:"min_price"`
	PriceStep   float64                `json:"price_step"`
	CreatedAt   *customTime.CustomTime `json:"created_at"`
}
