//TODO:
//
//1. Make depth of search a variable
//5. Make one map[symbol_name]Struct_of_all_data for every symbol <- entry point for API aggregation
//6. Refactor for bulk variable names/code and pointers usage and GOROUTINES
//10. Redo with websocket
//11. Handle API errors
//12. goroutine for counting opportunities
//14. refresh each iteration's price by pointers
//15. add revoking goroutine that returns unused coins back to basset - needs global lock and not trading state
//16. move api to separate file
//19. FOK mode - ?
//21. partiallyFilled error
//22. goroutinize: simultaniously trade and watch for another opportunity
//23. fix real % increase

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/astro/triarb/pkg/arbitrage"
	"github.com/astro/triarb/pkg/exchange"
	"github.com/astro/triarb/pkg/models"
	"github.com/astro/triarb/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Exchange struct {
		HitBTC struct {
			APIKey    string `yaml:"api_key"`
			SecretKey string `yaml:"secret_key"`
			BaseURL   string `yaml:"base_url"`
		} `yaml:"hitbtc"`
	} `yaml:"exchange"`
	Trading struct {
		MinProfitPercent          float64 `yaml:"min_profit_percent"`
		MinPairExistenceTime      int     `yaml:"min_pair_existence_time"`
		PriceAdjustmentMultiplier int     `yaml:"price_adjustment_multiplier"`
		TestingMode               bool    `yaml:"testing_mode"`
		VerboseMode               bool    `yaml:"verbose_mode"`
	} `yaml:"trading"`
}

func loadConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func loadCoins(path string) ([]string, error) {
	return utils.ReadLines(path)
}

func main() {
	var (
		configPath string
		amount     float64
		currency   string
		infinite   bool
		listen     bool
	)

	flag.StringVar(&configPath, "config", "config/config.yaml", "path to config file")
	flag.Float64Var(&amount, "a", 0.0, "amount of currency to trade")
	flag.StringVar(&currency, "c", "", "currency to trade, uppercase, abbreviation")
	flag.BoolVar(&infinite, "infinite", false, "loop over cur actual prices")
	flag.BoolVar(&listen, "listen", false, "wait 'til first opportunity, then exit")
	flag.Parse()

	if currency == "" {
		log.Fatal("currency (-c) is required")
	}

	// Load configuration
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Load supported coins
	coins, err := loadCoins("assets/coins.txt")
	if err != nil {
		log.Printf("Warning: could not load coins.txt: %v", err)
		coins = []string{} // Empty list if file not found
	}

	// Initialize exchange
	hitbtc := exchange.NewHitBTC(config.Exchange.HitBTC.APIKey, config.Exchange.HitBTC.SecretKey)

	// Initialize arbitrage engine
	engine := arbitrage.NewEngine(hitbtc, &arbitrage.Config{
		MinProfitPercent:          config.Trading.MinProfitPercent,
		MinPairExistenceTime:      config.Trading.MinPairExistenceTime,
		PriceAdjustmentMultiplier: config.Trading.PriceAdjustmentMultiplier,
		TestingMode:               config.Trading.TestingMode,
		VerboseMode:               config.Trading.VerboseMode,
		SupportedCoins:            coins,
	})

	// Check initial balance
	balance, err := hitbtc.GetBalance(currency)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	if !config.Trading.TestingMode && balance.Available < amount {
		log.Fatalf("Insufficient balance. Available: %f %s, Required: %f %s",
			balance.Available, currency, amount, currency)
	}

	fmt.Printf("Starting arbitrage with %f %s\n", amount, currency)

	// Main trading loop
	for {
		opportunities, err := engine.FindOpportunities(currency, amount)
		if err != nil {
			log.Printf("Error finding opportunities: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if len(opportunities) > 0 {
			fmt.Printf("Found %d opportunities\n", len(opportunities))

			// Group opportunities into trading sequences
			sequences := groupIntoSequences(opportunities)

			for _, sequence := range sequences {
				if err := engine.ExecuteTrade(sequence); err != nil {
					log.Printf("Error executing trade: %v", err)
					continue
				}

				if listen {
					return
				}
			}
		}

		if !infinite {
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}

// groupIntoSequences groups trading symbols into sequences of 3
func groupIntoSequences(symbols []models.TradingSymbol) [][]models.TradingSymbol {
	var sequences [][]models.TradingSymbol
	for i := 0; i < len(symbols); i += 3 {
		end := i + 3
		if end > len(symbols) {
			end = len(symbols)
		}
		sequences = append(sequences, symbols[i:end])
	}
	return sequences
}
