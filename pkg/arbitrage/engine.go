package arbitrage

import (
	"fmt"
	"time"

	"github.com/astro/triarb/pkg/exchange"
	"github.com/astro/triarb/pkg/models"
)

// Engine represents the arbitrage trading engine
type Engine struct {
	exchange *exchange.Exchange
	config   *Config
	symbols  []models.Symbol
	tickers  []models.Ticker
}

// Config represents arbitrage engine configuration
type Config struct {
	MinProfitPercent          float64
	MinPairExistenceTime      int
	PriceAdjustmentMultiplier int
	TestingMode               bool
	VerboseMode               bool
	SupportedCoins            []string
}

// NewEngine creates a new arbitrage engine
func NewEngine(exchange *exchange.Exchange, config *Config) *Engine {
	return &Engine{
		exchange: exchange,
		config:   config,
	}
}

// FindOpportunities searches for arbitrage opportunities
func (e *Engine) FindOpportunities(baseCurrency string, amount float64) ([]models.TradingSymbol, error) {
	// Get current market data
	symbols, err := e.exchange.GetSymbols()
	if err != nil {
		return nil, fmt.Errorf("failed to get symbols: %v", err)
	}
	e.symbols = symbols

	tickers, err := e.exchange.GetTickers()
	if err != nil {
		return nil, fmt.Errorf("failed to get tickers: %v", err)
	}
	e.tickers = tickers

	// Find trading sequences
	sequences := e.findTradingSequences(baseCurrency)

	// Calculate profits for each sequence
	profitableSequences := make([]models.TradingSymbol, 0)
	for _, sequence := range sequences {
		profit, tradingSymbols := e.calculateProfit(sequence, amount)
		if profit >= e.config.MinProfitPercent {
			profitableSequences = append(profitableSequences, tradingSymbols...)
		}
	}

	return profitableSequences, nil
}

// findTradingSequences finds all possible trading sequences
func (e *Engine) findTradingSequences(baseCurrency string) [][]models.Symbol {
	// Implementation of the trading sequence finding logic
	// This is where we'll move the existing tradingSequence function
	return nil
}

// calculateProfit calculates the potential profit for a trading sequence
func (e *Engine) calculateProfit(sequence []models.Symbol, amount float64) (float64, []models.TradingSymbol) {
	// Implementation of the profit calculation logic
	// This is where we'll move the existing percentCalculation function
	return 0, nil
}

// ExecuteTrade executes a trading sequence
func (e *Engine) ExecuteTrade(sequence []models.TradingSymbol) error {
	if e.config.TestingMode {
		return e.simulateTrade(sequence)
	}
	return e.executeRealTrade(sequence)
}

// simulateTrade simulates a trade without executing it
func (e *Engine) simulateTrade(sequence []models.TradingSymbol) error {
	for _, trade := range sequence {
		fmt.Printf("Simulated %s %f %s at %f\n",
			trade.Action, trade.Quantity, trade.Symbol.Id, trade.Price)
	}
	return nil
}

// executeRealTrade executes a real trade on the exchange
func (e *Engine) executeRealTrade(sequence []models.TradingSymbol) error {
	for _, trade := range sequence {
		order, err := e.exchange.PlaceOrder(trade.Symbol.Id, trade.Action, trade.Quantity, trade.Price)
		if err != nil {
			return fmt.Errorf("failed to place order: %v", err)
		}

		// Wait for order to be filled
		for order.Status != "filled" {
			time.Sleep(200 * time.Millisecond)
			order, err = e.exchange.GetOrder(order.ClientOrderId)
			if err != nil {
				return fmt.Errorf("failed to get order status: %v", err)
			}
		}
	}
	return nil
}
