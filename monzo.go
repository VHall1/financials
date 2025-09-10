package main

import (
	"context"
	"encoding/json"
	"net/http"
)

const monzoBaseURL = "https://api.monzo.com/"

type MonzoClient struct {
	baseURL string
	token   string
}

func NewMonzoClient(token string) *MonzoClient {
	return &MonzoClient{
		baseURL: monzoBaseURL,
		token:   token,
	}
}

func (c *MonzoClient) doRequest(ctx context.Context, method, endpoint string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// https://docs.monzo.com/#transactions
type MonzoTransaction struct {
	// The amount of the transaction in minor units of
	// currency. For example pennies in the case of GBP. A
	// negative amount indicates a debit (most card transactions
	// will have a negative amount)
	Amount      int64
	Created     string
	Currency    string
	Description string
	ID          string
	Merchant    string
	Notes       string
	// Top-ups to an account are represented as transactions
	// with a positive amount and is_load = true. Other
	// transactions such as refunds, reversals or chargebacks
	// may have a positive amount but is_load = false
	IsLoad bool
	// The timestamp at which the transaction settled. In most
	// cases, this happens 24-48 hours after created. If this
	// field is an empty string, the transaction is authorised but
	// not yet "complete."
	Settled string
	// The category can be set for each transaction by the user.
	// Over time we learn which merchant goes in which
	// category and auto-assign the category of a transaction. If
	// the user hasn't set a category, we'll return the default
	// category of the merchant on this transactions. Top-ups
	// have category mondo. Valid values are general,
	// eating_out, expenses, transport, cash, bills,
	// entertainment, shopping, holidays, groceries.
	Category string
}

func (c *MonzoClient) ListTransactions(ctx context.Context, accountID string) ([]MonzoTransaction, error) {
	type response struct {
		Transactions []MonzoTransaction
	}

	resp, err := c.doRequest(ctx, http.MethodGet, "transactions?account_id="+accountID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody response
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return respBody.Transactions, nil
}
