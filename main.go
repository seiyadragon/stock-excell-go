package main

import (
	"fmt"
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/equity"
	"github.com/shopspring/decimal"
)

type Stock struct {
	quote             finance.Equity
	growth5yrs        float64
	growth5yrsPercent float64
	potentialEarning  float64
	potentialLoss     float64
}

func newStock(symbol string) *Stock {
	q, err := equity.Get(symbol)
	if err != nil {
		return (nil)
	}

	params := &chart.Params{
		Symbol:   q.Quote.Symbol,
		Interval: datetime.OneMonth,
		Start:    &datetime.Datetime{Month: 1, Day: 1, Year: 2000},
		End:      &datetime.Datetime{Month: int(time.Now().Month()), Day: time.Now().Day(), Year: time.Now().Year()},
	}
	iter := chart.Get(params)

	var priceArray []decimal.Decimal
	//var growth5yrs float64 = 0
	//var growth5yrsPercent float64 = 0
	//var monthCounter int = 0

	for iter.Next() {
		priceArray = append(priceArray, iter.Bar().Close)
	}

	println(len(priceArray))

	if err := iter.Err(); err != nil {
		return nil
	}

	print("[" + q.Symbol + "] Price: " + fmt.Sprintf("%f", q.RegularMarketPrice))

	return &Stock{*q, 0, 0, 0, 0}
}

func main() {
	newStock("f")
}
