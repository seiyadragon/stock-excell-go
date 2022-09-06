package main

import (
	"log"
	"strconv"
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/equity"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

type Stock struct {
	quote             finance.Equity
	growth5yrs        float64
	growth5yrsPercent float64
	potentialEarning  float64
	potentialLoss     float64
	risk              float64
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
	var risk float64 = 0

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

	potentialEarning = growth5yrs + (q.TrailingAnnualDividendRate * 5)
	potentialLoss = q.Quote.RegularMarketPrice - (q.TrailingAnnualDividendYield * 5)

	risk = potentialLoss / potentialEarning

	if growth5yrs < 0 || growth5yrsPercent < 0.8 || q.ForwardPE > 25 {
		return nil
	}

	println("[" + q.Symbol + "] Loaded...")

	return &Stock{*q, growth5yrs, growth5yrsPercent, potentialEarning, potentialLoss, risk}
}

func main() {
	var stocks []Stock

	stocks = append(stocks, *newStock("abr"))
	stocks = append(stocks, *newStock("F"))
	stocks = append(stocks, *newStock("msft"))

	excel := excelize.NewFile()
	excel.SetCellValue("Sheet1", "A1", "Symbol")
	excel.SetCellValue("Sheet1", "B1", "Price")
	excel.SetCellValue("Sheet1", "C1", "Dividend")
	excel.SetCellValue("Sheet1", "D1", "Dividend 5yrs")
	excel.SetCellValue("Sheet1", "E1", "PER")
	excel.SetCellValue("Sheet1", "F1", "Growth 5yrs")
	excel.SetCellValue("Sheet1", "G1", "Time Growing")
	excel.SetCellValue("Sheet1", "H1", "Potential Income")
	excel.SetCellValue("Sheet1", "I1", "Potential Loss")
	excel.SetCellValue("Sheet1", "J1", "Risk")

	excel.SetColWidth("Sheet1", "A", "J", 20)

	cellAlphabet := "ABCDEFGHIJ"

	for i := 0; i < len(stocks); i++ {
		for _, c := range cellAlphabet {
			index := string(string(c) + strconv.FormatInt(int64(i+2), 10))

			if c == 'A' {
				excel.SetCellValue("Sheet1", index, stocks[i].quote.Symbol)
			}
			if c == 'B' {
				excel.SetCellValue("Sheet1", index, stocks[i].quote.Quote.RegularMarketPrice)
			}
			if c == 'C' {
				excel.SetCellValue("Sheet1", index, stocks[i].quote.TrailingAnnualDividendYield)
			}
			if c == 'D' {
				excel.SetCellValue("Sheet1", index, stocks[i].quote.TrailingAnnualDividendRate*5)
			}
			if c == 'E' {
				excel.SetCellValue("Sheet1", index, stocks[i].quote.ForwardPE)
			}
			if c == 'F' {
				excel.SetCellValue("Sheet1", index, stocks[i].growth5yrs)
			}
			if c == 'G' {
				excel.SetCellValue("Sheet1", index, stocks[i].growth5yrsPercent)
			}
			if c == 'H' {
				excel.SetCellValue("Sheet1", index, stocks[i].potentialEarning)
			}
			if c == 'I' {
				excel.SetCellValue("Sheet1", index, stocks[i].potentialLoss)
			}
			if c == 'J' {
				excel.SetCellValue("Sheet1", index, stocks[i].risk)
			}
		}
	}

	if err := excel.SaveAs("simple.xlsx"); err != nil {
		log.Fatal(err)
	}
}
