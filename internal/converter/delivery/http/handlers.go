package http

import (
	"encoding/json"
	"erpv1/internal/converter"
	"net/http"
)

type converterHandlers struct {
	usecase converter.UseCase
}

func NewConverterHandlers(usecase converter.UseCase) converter.Handlers {
	return &converterHandlers{usecase: usecase}
}

func (h *converterHandlers) GetExchangeRatesFromCB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var currentCountry string
		country, ok := r.URL.Query()["country"]
		if !ok {
			currentCountry = "ru"
		} else {
			currentCountry = country[0]
		}

		data, err := h.usecase.GetExchangeRatesFromCB(currentCountry)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(data.Rates)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Произошла внутренняя ошибка."))
			return
		}

		w.Write(res)
	}
}

func (h *converterHandlers) ConvertRate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var currentCountry string
		country, ok := r.URL.Query()["country"]
		if !ok {
			currentCountry = "ru"
		} else {
			currentCountry = country[0]
		}
		from, ok := r.URL.Query()["from"]
		if !ok || from[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("from is required"))
			return
		}
		to, ok := r.URL.Query()["to"]
		if !ok || to[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("to is required"))
			return
		}
		value, ok := r.URL.Query()["value"]
		if !ok || value[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("value is required"))
			return
		}

		data, err := h.usecase.ConvertRate(currentCountry, from[0], to[0], value[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Произошла ошибка!"))
			return
		}

		w.Write(res)
	}
}
