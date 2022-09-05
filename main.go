package main

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

func main() {
	var stockList [10]finance.Quote

	q, err := quote.Get("ABR")
	if err != nil {
		panic(err)
	}

	println(q.Symbol)
}
