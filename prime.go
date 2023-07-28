package coprime

import (
	"net/http"
	"time"
)

const (
    API_URL = "https://api.prime.coinbase.com"
    API_VERSION = "v1"

)
var methods = []string{
    "allocations",
    "getAlloctations",
}

type Coprime struct {
    BaseURL string
    Secret string
    Key string
    Passphrase string
    HTTPClient *http.Client
    RetryCount int
}

func New(secret, key, passphrase string) *Coprime {

    cp := Coprime {
        BaseURL: API_URL,
        Secret: secret,
        Key: key,
        Passphrase: passphrase,
        HTTPClient: &http.Client{
            Timeout: 15 * time.Second,
        },
        RetryCount: 0,

    }
    return &cp
}

func (c *Coprime) SetBaseURL(url string) {
    c.BaseURL = url
}

func (c *Coprime) SetRetryCount(retries int) {
    c.RetryCount = retries
}

fung (C *Coprime) getHeaders() {
    h: make(map[string]string)
    h["CB-ACCESS-KEY"] = c.Key
	h["CB-ACCESS-PASSPHRASE"] = c.Passphrase
	h["CB-ACCESS-TIMESTAMP"] = timestamp

    message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)


}




