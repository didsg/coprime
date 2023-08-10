package coprime

import (
	"errors"
	"testing"
)

func TestCreateLimitOrders(t *testing.T) {
	client := NewTestClient()

	order := ProOrder{
		Price:     "1.00000000",
		Size:      "1000.00000000",
		Side:      "buy",
		ProductID: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&order)
	if err != nil {
		t.Error(err)
	}

	if savedOrder.ID == "" {
		t.Error(errors.New("No create id found"))
	}

	props := []string{"Price", "Size", "Side", "ProductID"}
	_, err = CompareProperties(order, savedOrder, props)
	if err != nil {
		t.Error(err)
	}

	if err := client.CancelProOrder(savedOrder.ID); err != nil {
		t.Error(err)
	}
}

func TestCreateMarketOrders(t *testing.T) {
	client := NewTestClient()

	order := ProOrder{
		Funds:     "10.00",
		Size:      "1000.00000000",
		Side:      "buy",
		Type:      "market",
		ProductID: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&order)
	if err != nil {
		t.Error(err)
	}

	if savedOrder.ID == "" {
		t.Error(errors.New("No create id found"))
	}

	props := []string{"Price", "Size", "Side", "ProductID"}
	_, err = CompareProperties(order, savedOrder, props)
	if err != nil {
		t.Error(err)
	}
}

func TestCancelOrder(t *testing.T) {
	client := NewTestClient()

	order := ProOrder{
		Price:     "1.00",
		Size:      "1000.00",
		Side:      "buy",
		ProductID: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&order)
	if err != nil {
		t.Error(err)
	}

	if err := client.CancelProOrder(savedOrder.ID); err != nil {
		t.Error(err)
		t.Error(err)
	}
}

func TestGetOrder(t *testing.T) {
	client := NewTestClient()

	order := ProOrder{
		Price:     "1.00",
		Size:      "1.00",
		Side:      "buy",
		ProductID: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&order)
	if err != nil {
		t.Error(err)
	}

	getOrder, err := client.GetProOrder(savedOrder.ID)
	if err != nil {
		t.Error(err)
	}

	if getOrder.ID != savedOrder.ID {
		t.Error(errors.New("Order ids do not match"))
	}

	if err := client.CancelProOrder(savedOrder.ID); err != nil {
		t.Error(err)
	}
}

func TestListOrders(t *testing.T) {
	client := NewTestClient()
	cursor := client.ListProOrders()
	var orders []ProOrder

	for cursor.HasMore {
		if err := cursor.NextPage(&orders); err != nil {
			t.Error(err)
		}

		for _, o := range orders {
			if StructHasZeroValues(o) {
				t.Error(errors.New("Zero value"))
			}
		}
	}

	cursor = client.ListProOrders(ListProOrdersParams{Status: "open", ProductID: "BTC-EUR"})
	for cursor.HasMore {
		if err := cursor.NextPage(&orders); err != nil {
			t.Error(err)
		}

		for _, o := range orders {
			if StructHasZeroValues(o) {
				t.Error(errors.New("Zero value"))
			}
		}
	}
}

func TestCancelAllOrders(t *testing.T) {
	client := NewTestClient()

	for _, pair := range []string{"BTC-USD"} {
		for i := 0; i < 2; i++ {
			order := ProOrder{Price: "100000.00", Size: "1.00", Side: "sell", ProductID: pair}

			if _, err := client.CreateOrder(&order); err != nil {
				t.Error(err)
			}
		}

		orderIDs, err := client.CancelAllProOrders(CancelAllProOrdersParams{ProductID: pair})
		if err != nil {
			t.Error(err)
		}

		if len(orderIDs) != 2 {
			t.Error("Did not cancel all orders")
		}
	}
}
