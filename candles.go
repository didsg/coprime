package coprime

func (c *Client) GetCandles() error {
	var products PortfolioProducts
	hasNext := true

	baseRequestURL := fmt.Sprintf("/products/%s/products", productID)
	requestURL := baseRequestURL

	var allProducts []PrimeProducts
	for hasNext {
		_, err := c.Request("GET", "pro", requestURL, nil, &products)
		if err != nil {
			return nil, err
		}
		allProducts = append(allProducts, products.Products...)
		hasNext = products.Pagination.HasNext
		requestURL = fmt.Sprintf("%s?%s", baseRequestURL, Encode(products.Pagination))
	}
	return allProducts, nil
}

//     response = requests.get(
//         "https://"
//         + "api.exchange.coinbase.com"
//         + f"/products/{product_id}/candles?granularity={granularity}&start={eod_epoch_time}&end={eod_epoch_time}",
//         auth=auth,
//     )
//     candle_json = response.json()
