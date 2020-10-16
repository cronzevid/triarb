# triarb

This bot searches for triangular arbitrage opportunity at HITBTC exchange.
An example: given BTC as arg, recursively searches for pairs like:
```
BTC -> other_cur; other_cur -> second_cur; second_cur -> BTC
```
In case of success, calculations on current prices are performed. Can be tuned with configurable profit percentage.

Requirements: go

Build: go build -o bits main.go

Usage:
```
Usage ./triarb: -a amount -c currency
  -a float
    	amount of currency to trade
  -c string
    	currency to trade, uppercase, abbreviation
  -infinite
    	loop over cur actual prices
  -listen
    	wait 'til first opportunity, then exit
  -m int
    	times the minimun price is added/substracted from buy/sell to make valuable offer (default 1)
  -t int
    	times/0.2sec pair should exist to pass to trade stage (default 5)
  -testing
    	set false to initiate trade (default true)
  -v	verbose mode
```
Also, put your API keys to hitbtc.go

Provides test and verbose modes.
