package coprime

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
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

func TestGetPortfolios(t *testing.T) {
	client := getClient()
	portfolios, err := client.GetPortfolios()
	if err != nil {
		log.Fatalln(err)
	}
    fmt.Println(portfolios)
}

func TestTime(t *testing.T) {
	client := getClient()

	time, err := client.GetTime()
	if err != nil {
		log.Fatalln(err)
	}
    fmt.Println(time)
}


func TestGetProducts(t *testing.T) {
    client := getClient()
    portfolioID := getPortfolioID()
    products, err := client.GetAvailableProducts(portfolioID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(products)

}

func TestGetBook(t *testing.T) {
    client := getClient()

    products, err := client.GetAllProducts()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(products)
}
