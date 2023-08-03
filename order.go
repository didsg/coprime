package coprime

import (
	"fmt"
	"time"
)

type PostOrder struct {
	Side             string `json:"side,omitempty"`
	Type             string `json:"type,omitempty"`
	BaseQuantity     string `json:"base_quantity,omitempty"`
	LimitPrice       string `json:"limit_price,omitempty"`
	PortfolioID      string `json:"portfolio_id,omitempty"`
	ProductID        string `json:"product_id,omitempty"`
	ClientOrderID    string `json:"client_order_id,omitempty"`
	QuoteValue       string `json:"quote_value,omitempty"`
	StartTime        string `json:"start_time,omitempty"`
	ExpiryTime       string `json:"expiry_time,omitempty"`
	TimeInForce      string `json:"time_in_force,omitempty"`
	STPID            string `json:"stp_id,omitempty"`
	DisplayQuoteSize string `json:"display_quote_size,omitempty"`
	DisplayBaseSize  string `json:"display_base_size,omitempty"`
	IsRaiseExact     bool   `json:"is_raise_exact,omitempty"`
}

type OrderID struct {
	OrderID string `json:"order_id"`
}

type CancelOrderID struct {
	ID string `json:"id"`
}

type CancelAllOrdersParams struct {
	ProductID string
}

type ListOrdersParams struct {
	Status     string
	ProductID  string
	Pagination PaginationParams
}

func (c *Client) CreateOrder(newOrder *PostOrder, portfolioID string) (OrderID, error) {
	var orderID OrderID

	url := fmt.Sprintf("/v1/portfolios/%s/order", portfolioID)
	_, err := c.Request("POST", "prime", url, newOrder, &orderID)

	return orderID, err
}

func (c *Client) CancelOrder(portfolioID string, orderID string) (CancelOrderID, error) {
	var cancelOrderID CancelOrderID
	url := fmt.Sprintf("/v1/portfolios/%s/orders/%s/cancel", portfolioID, orderID)
	_, err := c.Request("POST", "prime", url, nil, &cancelOrderID)
	return cancelOrderID, err
}

func (c *Client) GetOrder(portfolioID string, orderID string) (Order, error) {
	var order GetOrder

	url := fmt.Sprintf("/v1/portfolios/%s/orders/%s", portfolioID, orderID)
	_, err := c.Request("GET", "prime", url, nil, &order)
	return order.Order, err
}

func (c *Client) GetOpenOrders(portfolioID string, pair string) ([]Order, error) {
	var openOrders OpenOrders
	hasNext := true
	baseRequestURL := fmt.Sprintf("/v1/portfolios/%s/open_orders?product_ids=%s", portfolioID, pair)
	requestURL := baseRequestURL
	var orders []Order
	for hasNext {
		_, err := c.Request("GET", "prime", requestURL, nil, &openOrders)
		if err != nil {
			return nil, err
		}
		orders = append(orders, openOrders.Orders...)
		hasNext = openOrders.Pagination.HasNext
		requestURL = fmt.Sprintf("%s?%s", baseRequestURL, Encode(openOrders.Pagination))
	}

	return orders, nil

}

type Order struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	PortfolioID        string    `json:"portfolio_id"`
	ProductID          string    `json:"product_id"`
	Side               string    `json:"side"`
	ClientOrderID      string    `json:"client_order_id"`
	Type               string    `json:"type"`
	BaseQuantity       string    `json:"base_quantity"`
	QuoteValue         string    `json:"quote_value"`
	LimitPrice         string    `json:"limit_price"`
	StartTime          time.Time `json:"start_time"`
	ExpiryTime         time.Time `json:"expiry_time"`
	Status             string    `json:"status"`
	TimeInForce        string    `json:"time_in_force"`
	CreatedAt          time.Time `json:"created_at"`
	FilledQuantity     string    `json:"filled_quantity"`
	FilledValue        string    `json:"filled_value"`
	AverageFilledPrice string    `json:"average_filled_price"`
	Commission         string    `json:"commission"`
	ExchangeFee        string    `json:"exchange_fee"`
}

type OpenOrders struct {
	Orders     []Order               `json:"orders"`
	Pagination PrimePaginationParams `json:"pagination"`
}

type GetOrder struct {
	Order Order `json:"order"`
}
