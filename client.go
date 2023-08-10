package coprime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type keys struct {
    key string
    passphrase string
    secret string
}

type Client struct {

	// prime
	PrimeURL  string
	PrimeKeys keys

	// pro
	ProURL  string
	ProKeys keys

	// adavanced
	AdvancedURL  string
	AdvancedKeys keys

    // sandbox
    SandboxURL string

	HTTPClient *http.Client
	RetryCount int
}

type ClientConfig struct {
	BaseURL    string
	Key        string
	Passphrase string
	Secret     string
}

const (
    PRIME_URL = "https://api.prime.coinbase.com"
    PRO_URL = "https://api.exchange.coinbase.com"
    ADVANCED_URL = "https://api.coinbase.com"
    SANDBOX_URL = "https://api-public.sandbox.exchange.coinbase.com"
)

func NewClient(primeKeys, proKeys, advancedKeys keys) *Client {
	client := Client{
		PrimeURL:     PRIME_URL,
		PrimeKeys:    primeKeys,
		ProURL:       PRO_URL,
		ProKeys:      proKeys,
		AdvancedURL:  ADVANCED_URL,
		AdvancedKeys: advancedKeys,
        SandboxURL: SANDBOX_URL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		RetryCount: 5,
	}
	return &client
}

func (c *Client) UpdateEndpoint(kind ApiType, endpoint string) error {
    switch kind {
        case Pro:
            c.ProURL = endpoint
        case Prime:
            c.PrimeURL = endpoint
        case Advanced:
            c.AdvancedURL = endpoint
        case Sandbox:
            c.SandboxURL = endpoint
        default:
            return errors.New("Invalid API Type")
    }
    return nil
}

func (c *Client) Request(method string, kind ApiType, url string, params, result interface{}) (res *http.Response, err error) {
	for i := 0; i < c.RetryCount+1; i++ {

		retryDuration := time.Duration((math.Pow(2, float64(i))-1)/2*1000) * time.Millisecond
		time.Sleep(retryDuration)

		res, err = c.request(method, kind, url, params, result)
		if res != nil && res.StatusCode == 429 {
			continue
		} else {
			break
		}
	}
	return res, err
}

func (c *Client) request(method string, kind ApiType, url string,
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
	switch kind {
    case Pro:
		fullURL = fmt.Sprintf("%s%s", c.ProURL, url)
    case Prime:
		fullURL = fmt.Sprintf("%s%s", c.PrimeURL, url)
    case Advanced:
		fullURL = fmt.Sprintf("%s%s", c.AdvancedURL, url)
    case Sandbox:
        fullURL = fmt.Sprintf("%s%s", c.SandboxURL, url)
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

	h, err := c.Headers(method, url, timestamp, string(data), kind)
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
		log.Fatal(err)
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
func (c *Client) Headers(method, url, timestamp, data string, kind ApiType) (map[string]string, error) {
    // create headers map
    h := make(map[string]string)
    var keySet keys
    var useSignatureOrSign string
	switch kind {
	case Pro:
        keySet = c.ProKeys
        h["CB-ACCESS-KEY"] = keySet.key
        h["CB-ACCESS-PASSPHRASE"] = keySet.passphrase
        h["CB-ACCESS-TIMESTAMP"] = timestamp
        useSignatureOrSign = "CB-ACCESS-SIGN"
    case Sandbox:
        keySet = c.ProKeys
        h["CB-ACCESS-KEY"] = keySet.key
        h["CB-ACCESS-PASSPHRASE"] = keySet.passphrase
        h["CB-ACCESS-TIMESTAMP"] = timestamp
        useSignatureOrSign = "CB-ACCESS-SIGN"
    case Prime:
        keySet = c.PrimeKeys
        h["X-CB-ACCESS-KEY"] = keySet.key
        h["X-CB-ACCESS-PASSPHRASE"] = keySet.passphrase
        h["X-CB-ACCESS-TIMESTAMP"] = timestamp
        useSignatureOrSign = "X-CB-ACCESS-SIGNATURE"
    case Advanced:
        keySet = c.AdvancedKeys
        useSignatureOrSign = "CB-ACCESS-SIGN"
	default:
		return nil, errors.New("Invalid api type, please use pro or prime")
	}

	// Cannot have any query parameters in url otherwise will get invalid api key
	url = strings.Split(url, "?")[0]

	message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)

	sig, err := generateSig(message, keySet.secret)
	if err != nil {
		return nil, err
	}

	h[useSignatureOrSign] = sig
	return h, nil
}
