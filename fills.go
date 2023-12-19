package coprime

import (
	"fmt"
)

type Fill struct {
	Commission     string `json:"commission"`
	FilledQuantity string `json:"filled_quantity"`
	FilledValue    string `json:"filled_value"`
	ID             string `json:"id"`
	OrderID        string `json:"order_id"`
	Price          string `json:"price"`
	ProductID      string `json:"product_id"`
	Side           string `json:"side"`
	Time           string `json:"time"`
	Venue          string `json:"venue"`
}

type OrderFills struct {
	Fills      []Fill
	Pagination PrimePaginationParams
}

func (c *Client) GetFills(orderID, portfolioID string) ([]Fill, error) {
	var oFills OrderFills
	hasNext := true
	baseRequestURL := fmt.Sprintf("/v1/portfolios/%s/orders/%s/fills", portfolioID, orderID)
	requestURL := baseRequestURL

	var fills []Fill
	for hasNext {
		_, err := c.Request("GET", "prime", requestURL, nil, &oFills)
		if err != nil {
			return nil, err
		}
		fills = append(fills, oFills.Fills...)
		hasNext = oFills.Pagination.HasNext
		requestURL = fmt.Sprintf("%s?%s", baseRequestURL, Encode(oFills.Pagination))
	}
	return fills, nil
}
