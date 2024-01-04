package http

import (
	"erpv1/internal/converter"
	"fmt"
	"net/http"
)

func MapConverterRoutes(prefix string, mux *http.ServeMux, h converter.Handlers) {
	mux.HandleFunc(fmt.Sprintf("/%v/exchange", prefix), h.GetExchangeRatesFromCB())
	mux.HandleFunc(fmt.Sprintf("/%v/convert", prefix), h.ConvertRate())
}
