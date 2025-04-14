package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/astro/triarb/pkg/models"
)

// Exchange represents the HITBTC exchange
type Exchange struct {
	APIKey    string
	SecretKey string
	Client    *http.Client
	BaseURL   string
}

// NewHitBTC creates a new HITBTC exchange instance
func NewHitBTC(apiKey, secretKey string) *Exchange {
	return &Exchange{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Client:    &http.Client{Timeout: 10 * time.Second},
		BaseURL:   "https://api.hitbtc.com/api/3",
	}
}

// GetSymbols retrieves all trading symbols
func (e *Exchange) GetSymbols() ([]models.Symbol, error) {
	resp, err := e.Client.Get(fmt.Sprintf("%s/public/symbol", e.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	if err := json.Unmarshal(body, &symbols); err != nil {
		return nil, err
	}

	return symbols, nil
}

// GetTickers retrieves current ticker data
func (e *Exchange) GetTickers() ([]models.Ticker, error) {
	resp, err := e.Client.Get(fmt.Sprintf("%s/public/ticker", e.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tickers []models.Ticker
	if err := json.Unmarshal(body, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}

// PlaceOrder places a new order
func (e *Exchange) PlaceOrder(symbol, side string, quantity, price float64) (*models.Order, error) {
	// Implementation of order placement
	return nil, nil
}

// GetOrder retrieves order information
func (e *Exchange) GetOrder(orderID string) (*models.Order, error) {
	// Implementation of order retrieval
	return nil, nil
}

// GetBalance retrieves account balance
func (e *Exchange) GetBalance(currency string) (*models.Balance, error) {
	// Implementation of balance retrieval
	return nil, nil
}
