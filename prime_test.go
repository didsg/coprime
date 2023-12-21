package coprime

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const (
	TIME_IN_FORCE = "GOOD_UNTIL_CANCELLED"
	TRADE_TIMEOUT = 120
)

func getClient() *Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("COINBASE_PRIME_KEY")
	pass := os.Getenv("COINBASE_PRIME_PASSPHRASE")
	secret := os.Getenv("COINBASE_PRIME_SECRET")

	primeURL := "https://api.prime.coinbase.com"
	proURL := "https://api.pro.coinbase.com"
	client := NewClient(primeURL, proURL, key, pass, secret)
	return client

}

func getPortfolioID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portfolioID := os.Getenv("PORTFOLIO_ID")
	return portfolioID
}

func getOrderID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	orderID := os.Getenv("ORDER_ID")
	return orderID
}

func addOrder(client *Client, pair string, direction string, orderType string, size string, args map[string]string) string {
	price := args["price"]
	clientOrderID := uuid.New().String()

	now := time.Now()
	expiryTime := now.Add(time.Duration(TRADE_TIMEOUT) * time.Second)
	sExpiryTime := expiryTime.Format("2006-01-02T15:04:05Z")

	order := PostOrder{
		ProductID:     pair,
		Side:          direction,
		ClientOrderID: clientOrderID,
		Type:          orderType,
		BaseQuantity:  size,
		LimitPrice:    price,
		TimeInForce:   TIME_IN_FORCE,
		ExpiryTime:    sExpiryTime,
	}
	orderID, err := client.CreateOrder(&order, getPortfolioID(), "prime")
	fmt.Println(fmt.Sprintf("Order %v (%s) submitted for pair %s and will expire at %s", clientOrderID, orderID, pair, sExpiryTime))
	if err != nil {
		log.Fatal(err)
	}
	return orderID.OrderID
}

func TestGetFills(t *testing.T) {
	client := getClient()
	portfolioID := getPortfolioID()
	orderID := getOrderID()

	fills, err := client.GetFills("prime", portfolioID, orderID)
	if err != nil {
		t.Error(err)
	}
	if len(fills) != 1 {
		t.Errorf("incorrect fill length")
	}
	for _, fill := range fills {
		log.Println(fill)
	}
}
func TestGetOpenOrders(t *testing.T) {
	client := getClient()
	portfolioID := getPortfolioID()
	pair := "BTC-USD"
	args := make(map[string]string, 1)
	args["price"] = "0.1"

	orderID := addOrder(client, pair, "BUY", "LIMIT", "0.00016", args)

	openOrders, err := client.GetOpenOrders(portfolioID, pair, "prime")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Open orders: %v", openOrders))

	cancelOrderID, err := client.CancelOrder(portfolioID, orderID, "prime")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Cancel order ID: %v", cancelOrderID))
}

func TestGetOrder(t *testing.T) {
	client := getClient()
	portfolioID := getPortfolioID()
	pair := "BTC-USD"
	args := make(map[string]string, 1)
	args["price"] = "0.1"

	orderID := addOrder(client, pair, "BUY", "LIMIT", "0.00016", args)

	order, err := client.GetOrder(portfolioID, orderID, "prime")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Found order: %v", order))

	cancelOrderID, err := client.CancelOrder(portfolioID, orderID, "prime")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Cancel order ID: %v", cancelOrderID))

}

func TestPlaceOffMarketOrderAndCancel(t *testing.T) {
	client := getClient()

	direction := "BUY"
	pair := "BTC-USD"
	tradeType := "LIMIT"
	portfolioID := getPortfolioID()

	var tradeSize string

	products, err := client.GetAvailableProducts(portfolioID, "prime")
	if err != nil {
		t.Error(err)
	}

	// Loop through the data and pull out one element
	var productID string
	for _, product := range products {
		if product.ID == pair {
			tradeSize = product.BaseMinSize
			productID = product.ID
			break
		}
	}

	orderBook, err := client.GetBook(pair, 1)
	if err != nil {
		t.Error(err)
	}

	// Get the current best side price and size for Leg 1.
	// Order placed at 10% of bid, but order size cannot be less than $1
	sSidePrice := orderBook.Bids[0].Price
	sidePrice, err := strconv.ParseFloat(sSidePrice, 64)
	if err != nil {
		t.Error(err)
	}
	sidePrice = sidePrice * 0.1
	if err != nil {
		t.Error(err)
	}
	sSidePrice = strconv.FormatFloat(sidePrice, 'f', 2, 64)

	// Place order
	args := make(map[string]string, 1)
	args["price"] = sSidePrice
	orderID := addOrder(client, productID, direction, tradeType, tradeSize, args)

	fmt.Println("Sleeping for 5 seconds ...")
	time.Sleep(5 * time.Second)

	cancelOrderID, err := client.CancelOrder(portfolioID, orderID, "prime")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Cancel order ID: %v", cancelOrderID))
}

func TestGetBalances(t *testing.T) {
	client := getClient()
	portfolioBalances, err := client.GetPortfolioBalances(getPortfolioID(), "prime")
	for _, balances := range portfolioBalances.Balances {
		if err != nil {
			t.Error(err)
		}
		if balances.Amount != "0" {
			fmt.Println(fmt.Sprintf("Currency: %s, Balance: %v", balances.Symbol, balances.Amount))
		}
	}
}

func TestGetPortfolios(t *testing.T) {
	client := getClient()
	portfolios, err := client.GetPortfolios("prime")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(portfolios)
}

func TestTime(t *testing.T) {
	client := getClient()

	time, err := client.GetTime()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(time)
}

func TestGetPortfolio(t *testing.T) {
	client := getClient()
	portfolioID := getPortfolioID()
	portfolio, err := client.GetPortfolio(portfolioID, "prime")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(portfolio)
}

func TestGetProducts(t *testing.T) {
	client := getClient()
	portfolioID := getPortfolioID()
	products, err := client.GetAvailableProducts(portfolioID, "prime")
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(products)
	// Loop through the data and pull out one element
	for _, product := range products {
		if product.ID == "BTC-USD" {
			fmt.Println(product.BaseMinSize)
			break
		}
	}
}

func TestGetBook(t *testing.T) {
	client := getClient()
	products, err := client.GetAllProducts()
	if err != nil {
		t.Error(err)
	}

	// Loop through the data and pull out one element
	for _, product := range products {
		if product.ID == "BTC-USD" {
			fmt.Println(product)

			book, err := client.GetBook(product.ID, 1)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(book)
			break
		}
	}
}
