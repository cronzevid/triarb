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

import ("encoding/json"
        "os"
        "flag"
        "time"
        "math"
        "bufio"
        "fmt"
        "triarb/toolbelt"
        "triarb/hitbtc")

type Ticker struct {
  Bid float64 `json:",string"`
  Ask float64 `json:",string"`
  Last float64 `json:",string"`
  Symbol string
}
type Symbol struct {
  Id string
  BaseCurrency string
  QuoteCurrency string
  FeeCurrency string
  TakeLiquidityRate float64 `json:",string"`
  ProvideLiquidityRate float64 `json:",string"`
  QuantityIncrement float64  `json:",string"`
  TickSize float64 `json:",string"`
}
type Last struct {
  Currency Symbol
  Action string
}
type SymbolMap struct {
  BaseCurrency string
  QuoteCurrency string
  FeeCurrency string
  Bid float64
  Ask float64
  Last float64
}
type TradingSymbol struct {
  Symbol Symbol
  Action string
  Quantity float64
  Price float64
}
type OrderId struct {
  Id int
  ClientOrderId string
  Symbol string
  Side string
  Status string
  Type string
  TimeInForce string
  Quantity float64 `json:",string"`
  Price float64 `json:",string"`
  CumQuantity float64 `json:",string"`
  PostOnly bool `json:",bool"`
  CreatedAt string
  UpdatedAt string
}

func tradingAction(percent float64, sequence [3]TradingSymbol, valueCounter map[[3]TradingSymbol]int, feeSymbolTake map[string]float64, basset string, balance float64, times int, listen bool, testing bool) {

    timeTrack := time.Now()
    var pretty_out string

    //defer toolbelt.TimeTrack(time.Now(), "do trade")

    for _, element := range sequence {
         pretty_out = pretty_out + element.Symbol.Id + "|" + element.Action + "->"
    }
    pretty_out = pretty_out[0:len(pretty_out)-2]

    if len(valueCounter) != 0 {
        if _, ok := valueCounter[sequence]; ok {
            valueCounter[sequence] += 1
        } else {
            valueCounter[sequence] = 1
        }
    } else {
        valueCounter[sequence] = 1
    }

    if valueCounter[sequence] >= times {
        var starting_am float64
        var ending_am float64
        var actualBalance float64

        actualBalance = hitbtc.BalanceRequest(basset).Available
        fmt.Println("Real balance:", actualBalance, basset)

        if testing == false {
            starting_am = actualBalance
        } else {
            if sequence[0].Action ==  "buy" {
                starting_am = sequence[0].Quantity * sequence[0].Price
            } else if sequence[0].Action ==  "sell" {
                starting_am = sequence[0].Quantity
            }
            if sequence[2].Action ==  "buy" {
                ending_am = sequence[2].Quantity
            } else if sequence[2].Action ==  "sell" {
                ending_am = sequence[2].Quantity * sequence[2].Price
            }
        }
        fmt.Printf(timeTrack.Format("2006-01-02 15:04:05.99999"))
        fmt.Printf("/%v/Predicted increase by %v%% (%v,%v)/Counting: %v\n", pretty_out, percent, balance, balance + balance*(percent/100), valueCounter[sequence])

        var order hitbtc.OrderId

        for _, element := range sequence {
            fee := feeSymbolTake[element.Symbol.Id]
            if element.Action == "buy" {
                if actualBalance >= element.Price * element.Quantity {
                    fmt.Println("ORDER_ACTION", element.Symbol.Id, element.Action, element.Quantity, element.Price)
                    if testing == false {
                        order = hitbtc.Order(element.Symbol.Id, element.Action, element.Quantity, element.Price)
                        fmt.Printf("ORDER_DATA id:%v symbol:%v side:%v status:%v quantity:%v price:%v cumquantity:%v\n", order.ClientOrderId, order.Symbol, order.Side, order.Status, order.Quantity, order.Price, order.CumQuantity)
                    } else {
                        actualBalance = element.Quantity
                    }
                    fmt.Printf("Balance: %v, Price: %v, Quant: %v, Fee: %v\n", actualBalance, element.Price, element.Quantity, fee)
                } else {
                    fmt.Println("Not enough balance")
                    fmt.Printf("Balance: %v, Price: %v, Quant: %v, Fee: %v\n", actualBalance, element.Price, element.Quantity, fee)
                    break
                }
            } else if element.Action == "sell" {
                if actualBalance >= element.Quantity {
                    fmt.Println("ORDER_ACTION", element.Symbol.Id, element.Action, element.Quantity, element.Price)
                    if testing == false {
                        order = hitbtc.Order(element.Symbol.Id, element.Action, element.Quantity, element.Price)
                        fmt.Printf("ORDER_DATA id:%v symbol:%v side:%v status:%v quantity:%v price:%v cumquantity:%v\n", order.ClientOrderId, order.Symbol, order.Side, order.Status, order.Quantity, order.Price, order.CumQuantity)
                    } else {
                        actualBalance = element.Quantity * element.Price
                    }
                    fmt.Printf("Balance: %v, Price: %v, Quant: %v, Fee: %v\n", actualBalance, element.Price, element.Quantity, fee)
                } else {
                    fmt.Println("Not enough balance")
                    fmt.Printf("Balance: %v, Price: %v, Quant: %v, Fee: %v\n", actualBalance, element.Price, element.Quantity, fee)
                    break
                }
            }
            if testing == false {
                //orderHistory := hitbtc.OrderHistory(order.ClientOrderId)
                orderHistory := order.Status
                fmt.Println("ORDER_STATUS", orderHistory)
                //dotCounter := 0
                for orderHistory != "filled" {
                    if orderHistory == "canceled" {
                        fmt.Printf("FATAL: Order %v canceled, need manual interaction\n", order.ClientOrderId)
                        os.Exit(1)
                    }

                    //toolbelt.DotCounter(0)

                    time.Sleep(200 * time.Millisecond)
                    orderHistory = hitbtc.OrderHistory(order.ClientOrderId).Status
                }
                if element.Action == "buy" {
                    actualBalance = hitbtc.BalanceRequest(element.Symbol.BaseCurrency).Available
                } else if element.Action == "sell" {
                    actualBalance = hitbtc.BalanceRequest(element.Symbol.QuoteCurrency).Available
                }
            }
        }
        fmt.Printf(timeTrack.Format("2006-01-02 15:04:05.99999"))
        if testing == false {
            fmt.Printf("/%v/Real balance increase by %v%% (%v,%v)\n", pretty_out, (actualBalance/starting_am - 1)*100, starting_am, actualBalance)
        } else {
            fmt.Printf("/%v/Real balance increase by %v%% (%v,%v)\n", pretty_out, (ending_am/starting_am - 1)*100, starting_am, ending_am)
        }
        fmt.Println("-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_\n")
        if listen {
            os.Exit(0)
        }
    }

    return
}

func tradingSequence(all_symbols_slice []Symbol, algo []map[Symbol]string, balance_coin string, intermid_coin string) []map[Symbol]string {

    //
    //profiling
    //defer toolbelt.TimeTrack(time.Now(), "sequence")
    //
    var res []map[Symbol]string
    var last Last

    //
    //there may be no similar values in row
    //
    for _, element := range all_symbols_slice {
        if len(algo) > 0 {
            last_val := algo[len(algo)-1]
            for k, _ := range last_val {
              last = Last{Currency: k}
            }
        } else {
            last = Last{Action: ""}
        }

        //
        //two cases for 'sell' and 'buy' assignment
        //
        if element.QuoteCurrency == intermid_coin && element != last.Currency {
            if element.BaseCurrency != balance_coin {
                //
                //3 is depth of value search (btcusd - usdxmr  - xmrbtc)
                //
                if len(algo) == 3 {
                    break
                }
                //
                //if no balance coin in buy type, go recursive, else end if
                //
                //fmt.Println("Buy", element.Id, "for", element.QuoteCurrency)
                algo = append(algo, map[Symbol]string{element: "buy"})
                res = append(res, tradingSequence(all_symbols_slice, algo, balance_coin, element.BaseCurrency)...)
                algo = algo[0:len(algo)-1]
            } else if element.BaseCurrency == balance_coin {
                //fmt.Println("Buy", element.Id, "for", element.QuoteCurrency)
                algo = append(algo, map[Symbol]string{element: "buy"})
                //fmt.Println("/---Stop---/")
                return algo[len(algo)-3:len(algo)]
            }
        } else if element.BaseCurrency == intermid_coin && element != last.Currency {
            if element.QuoteCurrency != balance_coin {
                if len(algo) == 3 {
                    break
                }
                //fmt.Println("Sell", element.Id, "for", element.BaseCurrency)
                algo = append(algo, map[Symbol]string{element: "sell"})
                res = append(res, tradingSequence(all_symbols_slice, algo, balance_coin, element.QuoteCurrency)...)
                algo = algo[0:len(algo)-1]
            } else if element.QuoteCurrency == balance_coin {
                //fmt.Println("Sell", element.Id, "for", element.BaseCurrency)
                algo = append(algo, map[Symbol]string{element: "sell"})
                //fmt.Println("/---Stop---/")
                return algo[len(algo)-3:len(algo)]
            }
        }
    }
    return res
}

func percentCalculation (sequence []map[Symbol]string, askMap map[string]float64, bidMap map[string]float64, feeSymbolTake map[string]float64, quantDivider map[string]float64, priceDivider map[string]float64, balance float64, mult int) map[float64][3]TradingSymbol {
    var percent float64
    var amountBuy float64
    var amountSell float64
    var sequence_list [3]TradingSymbol
    have := balance
    multi := float64(mult)

    //defer toolbelt.TimeTrack(time.Now(), "calculations")

    for i, currency := range sequence {
        //
        //Calculate only decreasing fee (not //feeProvide)
        //
        for symbol, act := range currency {
            feeTake := feeSymbolTake[symbol.Id]
            quat := quantDivider[symbol.Id]
            pr := priceDivider[symbol.Id]
            //substract fee, so to use less startin balance to control future fees
            have = have - have * feeTake
            if act == "buy" {
                //fmt.Printf("       | Before Buy | Each: %.15f Have: %.15f amountBuy: %.15f price: %.15f QuantDivider: %.15f PriceDivider: %.15f\n", askMap[symbol.Id], have, toolbelt.Round(amountBuy,quat), toolbelt.Round(have,pr), quat, pr)
                //count amount and round it
                amountBuy = have / askMap[symbol.Id]
                amountBuy = toolbelt.Round(amountBuy,quat)
                //count price from what we have and round it
                have = toolbelt.Round(have,pr)
                //
                //Detailed buy
                //
                //fmt.Printf("%v | Buy %.10f %v each %.10f %v\n", symbol.Id, amountBuy, symbol.BaseCurrency, askMap[symbol.Id], symbol.QuoteCurrency)
                //fmt.Printf("       | Base: %v Quote: %v Each: %.10f Have: %.10f Recieve: %.10f Total price: %.10f Total price(fee'd): %.10f\n", symbol.BaseCurrency, symbol.QuoteCurrency, askMap[symbol.Id], have, amountBuy, amountBuy * askMap[symbol.Id], amountBuy * askMap[symbol.Id] * (1 + feeTake))
                //fmt.Printf("       | Pricing %v instead of %v\n", askMap[symbol.Id] + multi*pr, askMap[symbol.Id])
                sequence_list[i] = TradingSymbol{symbol, "buy", amountBuy, askMap[symbol.Id] + multi*pr} //add mult times floorin for more attractive price
                have = amountBuy
            } else if act == "sell"{
                //fmt.Printf("       | Before Sell | Each: %.10f Have: %.10f amountSell: %.10f price: %.10f QuantDivider: %.10f PriceDivider: %.10f\n", bidMap[symbol.Id], have, toolbelt.Round(amountSell,quat), toolbelt.Round(have,pr), quat, pr)
                //round amount to sell
                amountSell = toolbelt.Round(have,quat)
                //count price from what we have and round it
                have = have * bidMap[symbol.Id]
                have = toolbelt.Round(have,pr)
                //
                //Detailed sell
                //
                //fmt.Printf("%v | Sell %.10f %v each %.10f %v\n", symbol.Id, amountSell, symbol.BaseCurrency, bidMap[symbol.Id], symbol.QuoteCurrency)
                //fmt.Printf("       | Base: %v Quote: %v Each: %.10f Have: %.10f Recieve: %.10f Total price: %.10f Total price(fee'd): %.10f\n", symbol.BaseCurrency, symbol.QuoteCurrency, bidMap[symbol.Id], have, amountSell * bidMap[symbol.Id], amountSell, amountSell - amountSell * feeTake)
                //fmt.Printf("       | Pricing %v instead of %v\n", bidMap[symbol.Id] - multi*pr, bidMap[symbol.Id])
                sequence_list[i] = TradingSymbol{symbol, "sell", amountSell, bidMap[symbol.Id] - multi*pr} //substract mult times floorin for more attractive price
            }
        }
    }
    fin_balance := (balance-have)/balance
    if math.Signbit(fin_balance) {
        abs := math.Abs(fin_balance) * 100
        percent = toolbelt.Round(abs,0.00125)
    }
    //fmt.Printf("Final balance: %v-%v Percent: %v\n\n", balance, have, percent)
    return map[float64][3]TradingSymbol{percent: sequence_list}
}

func main() {
    var ticker_json []Ticker
    var symbol_json []Symbol
    var clean_symbol_json []Symbol

    var balance float64
    flag.Float64Var(&balance, "a", 0, "amount of currency to trade")
    var basset string
    flag.StringVar(&basset, "c", "", "currency to trade, uppercase, abbreviation")

    var times int
    flag.IntVar(&times, "t", 5, "times/0.2sec pair should exist to pass to trade stage")
    var multiply int
    flag.IntVar(&multiply, "m", 1, "times the minimun price is added/substracted from buy/sell to make valuable offer")
    var verbose bool
    flag.BoolVar(&verbose, "v", false, "verbose mode")
    var testing bool
    flag.BoolVar(&testing, "testing", true, "set false to initiate trade")
    var infinite bool
    flag.BoolVar(&infinite, "infinite", false, "loop over cur actual prices")
    var listen bool
    flag.BoolVar(&listen, "listen", false, "wait 'til first opportunity, then exit")

    flag.Parse()

    if basset == "" || balance == 0 {
        fmt.Printf("Usage %v: -a amount -c currency\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }

    //defer toolbelt.TimeTrack(time.Now(), "main")

    fmt.Println("Starting...")

    symbols := hitbtc.DataRequest("https://api.hitbtc.com/api/2/public/symbol/")
    json.Unmarshal(symbols, &symbol_json)
    tickers := hitbtc.DataRequest("https://api.hitbtc.com/api/2/public/ticker/")
    json.Unmarshal(tickers, &ticker_json)

    if _, err := os.Stat("coins.txt"); err == nil {
        var coin_list []string
        top_coin_list := []string{"BTC", "ETH", "USD"}

        file, _ := os.Open("coins.txt")
        scanner := bufio.NewScanner(file)

        for scanner.Scan() {
            coin_list = append(coin_list, scanner.Text())
        }
        file.Close()

        //limit trades to coins.txt list
        fmt.Println("Popping symbols not in coins.txt ...")
        for _, coins := range symbol_json {
            if toolbelt.Contains(top_coin_list, coins.BaseCurrency) || toolbelt.Contains(top_coin_list, coins.QuoteCurrency) {
                if toolbelt.Contains(coin_list, coins.BaseCurrency) &&  toolbelt.Contains(coin_list, coins.QuoteCurrency) {
                    clean_symbol_json = append(clean_symbol_json, coins)
                }
            } else {
                if toolbelt.Contains(coin_list, coins.BaseCurrency) ||  toolbelt.Contains(coin_list, coins.QuoteCurrency) {
                    clean_symbol_json = append(clean_symbol_json, coins)
                }

            }
        }
    } else {
        fmt.Println("No coins.txt here, skipping popping ...")
        clean_symbol_json = symbol_json
    }

    bidMap := make(map[string]float64)
    askMap := make(map[string]float64)
    for _, prices := range ticker_json {
        bidMap[prices.Symbol] = prices.Bid
        askMap[prices.Symbol] = prices.Ask
    }

    feeSymbolTake := make(map[string]float64)
    feeSymbolProvide := make(map[string]float64)
    quantDivider := make(map[string]float64)
    priceDivider := make(map[string]float64)
    for _, i := range clean_symbol_json {
        feeSymbolTake[i.Id] = i.TakeLiquidityRate
        feeSymbolProvide[i.Id] = i.ProvideLiquidityRate
        quantDivider[i.Id] = i.QuantityIncrement
        priceDivider[i.Id] = i.TickSize
    }

    //Create full map with all data instead of several maps
    //fullMap := make(map[string]*SymbolMap)
    //for _, j := range symbol_json {
    //    fullMap[j.Id] = &SymbolMap{BaseCurrency: j.BaseCurrency, QuoteCurrency: j.QuoteCurrency, FeeCurrency: j.FeeCurrency}
    //}
    //for _, i := range ticker_json {
    //    fullMap[i.Symbol].Bid = i.Bid
    //}

    //get recursive func res
    var algo []map[Symbol]string
    res := tradingSequence(clean_symbol_json, algo, basset, basset)

    //split result by 3-size chunks
    var divided [][]map[Symbol]string
    chunkSize := 3

    for i := 0; i < len(res); i += chunkSize {
        end := i + chunkSize
        if end > len(res) {
            end = len(res)
        }
        divided = append(divided, res[i:end])
    }

    fmt.Println("Trading", balance, basset, "...\n")

    actualBalance := hitbtc.BalanceRequest(basset).Available
    if testing == false && actualBalance < balance {
        fmt.Println("Not enough balance:", actualBalance, balance)
        os.Exit(1)
    }


    valueCounter := make(map[[3]TradingSymbol]int)
    valueCounterReflect := make(map[[3]TradingSymbol]int)

    //infinite if cli arg
    for cycle := 0; cycle < times; cycle = cycle {
        percentSequenceRes := make(map[float64][3]TradingSymbol)

        // goroutine since we don't care about order
        for _, sequence := range divided {
            percentSequence := percentCalculation(sequence, askMap, bidMap, feeSymbolTake, quantDivider, priceDivider, balance, multiply)
            for percent, sequence := range percentSequence {
                if percent >= 0.1 {
                    percentSequenceRes[percent] = sequence
                }
            }
        }

        ///////////////
        var temp_ticker []Ticker
        tick := hitbtc.DataRequest("https://api.hitbtc.com/api/2/public/ticker/")
        json.Unmarshal(tick, &temp_ticker)
        bM := make(map[string]float64)
        aM := make(map[string]float64)
        for _, prices := range temp_ticker{
            bM[prices.Symbol] = prices.Bid
            aM[prices.Symbol] = prices.Ask
        }


        //set goroutine for each (not only best) result, only then do trade on persistent sequence
        if len(percentSequenceRes) != 0 {
            for percent, sequence := range percentSequenceRes {
                //for _, symbol := range sequence {
                //    fmt.Printf("Buy: %v, Sell: %v\n", bM[symbol.Symbol.Id], aM[symbol.Symbol.Id])
                //}
                tradingAction(percent, sequence, valueCounter, feeSymbolTake, basset, balance, times, listen, testing)
            }
        }

        //empty counters if nothing happens
        if len(valueCounterReflect) != 0 {
            for seq, count := range valueCounterReflect {
                if _, ok := valueCounter[seq]; ok {
                    if valueCounter[seq] != count + 1 {
                        delete(valueCounter, seq)
                        delete(valueCounterReflect, seq)
                    }
                } else {
                    delete(valueCounterReflect, seq)
                }
            }
        }
        //copy map 
        for k,v := range valueCounter {
            valueCounterReflect[k] = v
        }

        time.Sleep(200 * time.Millisecond)
        tickers = hitbtc.DataRequest("https://api.hitbtc.com/api/2/public/ticker/")
        json.Unmarshal(tickers, &ticker_json)
        for _, new_prices := range ticker_json {
            askMap[new_prices.Symbol] = new_prices.Ask
            bidMap[new_prices.Symbol] = new_prices.Bid
        }

        if infinite || listen {
            cycle = 0
        } else {
            cycle += 1
        }

    }
}
