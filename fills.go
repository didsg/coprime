package coprime

import "fmt"

//	func (c *Client) GetPortfolio(portfolioID string) (Portfolio, error) {
//		// var portfolio Portfolio
//		var portfolio Portfolio
//		requestURL := fmt.Sprintf("/v1/portfolios/%s", portfolioID)
//
//		_, err := c.Request("GET", "prime", requestURL, nil, &portfolio)
//		return portfolio, err
//	}
type Fill struct {
}

func (c *Client) GetFills(orderID, portfolioID string) (OrderID, error) {

	url := fmt.Sprintf("/v1/portfolios/%s/orders/%s/fills", portfolioID, orderID)
	_, err := c.Request("POST", "prime", url, newOrder, &orderID)

	// return orderID, err
}

// {'fills':
//      [
//          {'id': '3c337809-3af2-4891-a069-c3546c74c4ca', 'order_id': '5ddcc19f-a238-44ca-86ee-eff513fcceab', 'product_id': 'BTC-USD', 'side': 'SELL', 'filled_quantity': '0.00238697', 'filled_value': '99.42507517875701', 'price': '41653.257133', 'time': '2023-12-18T20:13:11.837Z', 'commission': '0.0497125375893785', 'venue': 'OTC'}
//      ],
//          'pagination':
//       {'next_cursor': '', 'sort_direction': 'DESC', 'has_next': False}
// }
