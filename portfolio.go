package coprime

import (
	"fmt"
	"time"
)

type Portfolio struct {
	Portfolio struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Entity       string `json:"entity_id"`
		Organization string `json:"organization_id"`
	} `json:"portfolio"`
}

type Portfolios struct {
	Portfolios []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		EntityID       string `json:"entity_id"`
		OrganizationID string `json:"organization_id"`
	} `json:"portfolios"`
}

type Account struct {
	ID        string `json:"id"`
	Balance   string `json:"balance"`
	Hold      string `json:"hold"`
	Available string `json:"available"`
	Currency  string `json:"currency"`
}

// Ledger

type LedgerEntry struct {
	ID        string        `json:"id,number"`
	CreatedAt time.Time     `json:"created_at,string"`
	Amount    string        `json:"amount"`
	Balance   string        `json:"balance"`
	Type      string        `json:"type"`
	Details   LedgerDetails `json:"details"`
}

type LedgerDetails struct {
	OrderID   string `json:"order_id"`
	TradeID   string `json:"trade_id"`
	ProductID string `json:"product_id"`
}

type GetAccountLedgerParams struct {
	Pagination PaginationParams
}

// Holds

type Hold struct {
	AccountID string    `json:"account_id"`
	CreatedAt time.Time `json:"created_at,string"`
	UpdatedAt time.Time `json:"updated_at,string"`
	Amount    string    `json:"amount"`
	Type      string    `json:"type"`
	Ref       string    `json:"ref"`
}

type ListHoldsParams struct {
	Pagination PaginationParams
}

func (c *Client) GetPortfolios() (Portfolios, error) {
	var portfolios Portfolios
	_, err := c.Request("GET", "prime", "/v1/portfolios", nil, &portfolios)
	return portfolios, err

}

func (c *Client) GetPortfolio(portfolioID string) (Portfolio, error) {
	// var portfolio Portfolio
	var portfolio Portfolio
	requestURL := fmt.Sprintf("/v1/portfolios/%s", portfolioID)

	_, err := c.Request("GET", "prime", requestURL, nil, &portfolio)
	return portfolio, err
}

func (c *Client) GetPortfolioBalances(portfolio_id string, currency ...string) (PortfolioBalances, error) {

	var portfolio_balances PortfolioBalances
	ccy := ""
	if len(currency) > 0 {
		ccy = fmt.Sprintf("&symbols=%s", currency[0])
	}
	url := fmt.Sprintf("/v1/portfolios/%s/balances?balance_type=TRADING_BALANCES%s", portfolio_id, ccy)
	_, err := c.Request("GET", "prime", url, nil, &portfolio_balances)
	return portfolio_balances, err

}

func (c *Client) ListAccountLedger(id string,

	p ...GetAccountLedgerParams) *Cursor {
	paginationParams := PaginationParams{}
	if len(p) > 0 {
		paginationParams = p[0].Pagination
	}
	return NewCursor(c, "GET", fmt.Sprintf("/accounts/%s/ledger", id), &paginationParams)

}

func (c *Client) ListHolds(id string, p ...ListHoldsParams) *Cursor {

	paginationParams := PaginationParams{}
	if len(p) > 0 {
		paginationParams = p[0].Pagination
	}
	return NewCursor(c, "GET", fmt.Sprintf("/accounts/%s/holds", id), &paginationParams)

}

type Balances struct {
	Symbol               string `json:"symbol"`
	Amount               string `json:"amount"`
	Holds                string `json:"holds"`
	BondedAmount         string `json:"bonded_amount"`
	ReservedAmount       string `json:"reserved_amount"`
	UnbondingAmount      string `json:"unbonding_amount"`
	UnvestedAmount       string `json:"unvested_amount"`
	PendingRewardsAmount string `json:"pending_rewards_amount"`
	PastRewardsAmount    string `json:"past_rewards_amount"`
	BondableAmount       string `json:"bondable_amount"`
	WithdrawableAmount   string `json:"withdrawable_amount"`
}

type TradingBalances struct {
	Total string `json:"total"`
	Holds string `json:"holds"`
}

type VaultBalances struct {
	Total string `json:"total"`
	Holds string `json:"holds"`
}

type PortfolioBalances struct {
	Balances        []Balances      `json:"balances"`
	Type            string          `json:"type"`
	TradingBalances TradingBalances `json:"trading_balances"`
	VaultBalances   VaultBalances   `json:"vault_balances"`
}
