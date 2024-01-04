package converter

import "net/http"

type Handlers interface {
	GetExchangeRatesFromCB() http.HandlerFunc
	ConvertRate() http.HandlerFunc
}
