package main

//APIResponse does things
type APIResponse struct {
	QuoteSummary struct {
		Result []struct {
			SummaryDetail struct {
				PreviousClose struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"previousClose"`
				DayLow struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"dayLow"`
				DayHigh struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"dayHigh"`

				Bid struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"bid"`
				Ask struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"ask"`
				Volume struct {
					Raw     int    `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"volume"`
				Currency string `json:"currency"`
			} `json:"summaryDetail"`
			Price struct {
				RegularMarketChangePercent struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketChangePercent"`
				RegularMarketChange struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketChange"`

				Exchange              string      `json:"exchange"`
				ExchangeName          string      `json:"exchangeName"`
				ExchangeDataDelayedBy int         `json:"exchangeDataDelayedBy"`
				MarketState           string      `json:"marketState"`
				QuoteType             string      `json:"quoteType"`
				Symbol                string      `json:"symbol"`
				UnderlyingSymbol      interface{} `json:"underlyingSymbol"`
				ShortName             string      `json:"shortName"`
				LongName              string      `json:"longName"`
				Currency              string      `json:"currency"`
				QuoteSourceName       string      `json:"quoteSourceName"`
				CurrencySymbol        string      `json:"currencySymbol"`
				FromCurrency          interface{} `json:"fromCurrency"`
				ToCurrency            interface{} `json:"toCurrency"`
				LastMarket            interface{} `json:"lastMarket"`
			} `json:"price"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteSummary"`
}
