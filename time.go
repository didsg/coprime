package coprime

import "fmt"

func (c *Client) GetTime() (ServerTime, error) {
	var serverTime ServerTime

	url := fmt.Sprintf("/time")
	_, err := c.Request("GET", "pro", url, nil, &serverTime)
	return serverTime, err
}

