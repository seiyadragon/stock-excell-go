package main

import (
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
	var priceDifArray []float64
	var priceDifPerArray []int

	var growth5yrs float64 = 0
	var growth5yrsPercent float64 = 0
	var potentialEarning float64 = 0
	var potentialLoss float64 = 0

	for iter.Next() {
		priceArray = append(priceArray, iter.Bar().Close)
	}

	if err := iter.Err(); err != nil {
		return nil
	}

	for i, j := 0, len(priceArray)-1; i < j; i, j = i+1, j-1 {
		priceArray[i], priceArray[j] = priceArray[j], priceArray[i]
	}

	for i := 0; i < 120; i++ {
		var one, _ = priceArray[i].Float64()
		var two, _ = priceArray[i+60].Float64()
		var result = one - two

		if result > 0 {
			priceDifPerArray = append(priceDifPerArray, 1)
		} else {
			priceDifPerArray = append(priceDifPerArray, 0)
		}

		priceDifArray = append(priceDifArray, result)
	}

	for i := 0; i < len(priceDifArray); i++ {
		growth5yrs += priceDifArray[i]
		growth5yrsPercent += float64(priceDifPerArray[i])
	}

	growth5yrs /= float64(len(priceDifArray))
	growth5yrsPercent /= float64(len(priceDifPerArray))

	potentialEarning = growth5yrs + (q.TrailingAnnualDividendYield * 5)
	potentialLoss = q.Quote.RegularMarketPrice - (q.TrailingAnnualDividendYield * 5)

	println("[" + q.Symbol + "] Loaded...")

	return &Stock{*q, growth5yrs, growth5yrsPercent, potentialEarning, potentialLoss}
}

func main() {
	newStock("abr")
}
