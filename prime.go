package coprime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	ProURL     string
	PrimeURL   string
	Secret     string
	Key        string
	Passphrase string
	HTTPClient *http.Client
	RetryCount int
}

type ClientConfig struct {
	BaseURL    string
	Key        string
	Passphrase string
	Secret     string
}

func NewClient(primeURL, proURL, primeKey, primePass, primeSecret string) *Client {

	client := Client{
		PrimeURL:   primeURL,
		ProURL:     proURL,
		Key:        primeKey,
		Passphrase: primePass,
		Secret:     primeSecret,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		RetryCount: 5,
	}

	return &client
}

func (c *Client) Request(method string, apiType string, url string, params, result interface{}) (res *http.Response, err error) {
	for i := 0; i < c.RetryCount+1; i++ {
		retryDuration := time.Duration((math.Pow(2, float64(i))-1)/2*1000) * time.Millisecond
		time.Sleep(retryDuration)
		res, err = c.request(method, apiType, url, params, result)
		if res != nil && res.StatusCode == 429 {
			continue
		} else {
			break
		}
	}
	return res, err
}

func (c *Client) request(method string, apiType string, url string,
	params, result interface{}) (res *http.Response, err error) {

	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, err
		}
		body = bytes.NewReader(data)
	}
	var fullURL string
	switch apiType {
	case "pro":
		fullURL = fmt.Sprintf("%s%s", c.ProURL, url)
	case "prime":
		fullURL = fmt.Sprintf("%s%s", c.PrimeURL, url)
	default:
		return nil, errors.New("Invalid api type, please use pro or prime")
	}
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	h, err := c.Headers(method, url, timestamp, string(data), apiType)
	if err != nil {
		return res, err
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = c.HTTPClient.Do(req)
	if err != nil {
		return res, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return res, err
	}
	_ = res.Body.Close()
	if res.StatusCode != 200 {
		coinbaseError := Error{}
		if err := json.Unmarshal(bodyBytes, &coinbaseError); err != nil {
			return res, err
		}
		return res, fmt.Errorf("%v", coinbaseError)
	}

	if result != nil {
		if err = json.Unmarshal(bodyBytes, result); err != nil {
			return res, err
		}
	}

	return res, nil
}

// Headers generates a map that can be used as headers to authenticate a request
func (c *Client) Headers(method, url, timestamp, data, apiType string) (map[string]string, error) {
	h := make(map[string]string)

	switch apiType {
	case "pro":
		h["CB-ACCESS-KEY"] = c.Key
		h["CB-ACCESS-PASSPHRASE"] = c.Passphrase
		h["CB-ACCESS-TIMESTAMP"] = timestamp

		// Cannot have any query parameters in url otherwise will get invalid api key
		url = strings.Split(url, "?")[0]

		message := fmt.Sprintf(
			"%s%s%s%s",
			timestamp,
			method,
			url,
			data,
		)

		sig, err := generateSig(message, c.Secret)
		if err != nil {
			return nil, err
		}
		h["CB-ACCESS-SIGN"] = sig
	case "prime":
		h["X-CB-ACCESS-KEY"] = c.Key
		h["X-CB-ACCESS-PASSPHRASE"] = c.Passphrase
		h["X-CB-ACCESS-TIMESTAMP"] = timestamp

		// Cannot have any query parameters in url otherwise will get invalid api key
		url = strings.Split(url, "?")[0]

		message := fmt.Sprintf(
			"%s%s%s%s",
			timestamp,
			method,
			url,
			data,
		)

		sig, err := generateSig(message, c.Secret)
		if err != nil {
			return nil, err
		}
		h["X-CB-ACCESS-SIGNATURE"] = sig
	default:
		return nil, errors.New("Invalid api type, please use pro or prime")
	}
	return h, nil
}

func (c *Client) GetTime() (ServerTime, error) {
	var serverTime ServerTime

	url := fmt.Sprintf("/time")
	_, err := c.Request("GET", "pro", url, nil, &serverTime)
	return serverTime, err
}

type ServerTime struct {
	ISO   string  `json:"iso"`
	Epoch float64 `json:"epoch,number"`
}

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
