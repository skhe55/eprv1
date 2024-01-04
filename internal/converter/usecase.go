package converter

import "erpv1/internal/models"

type UseCase interface {
	UpdateExchangeRatesFromRUCB() (models.CachedRates, error)
	GetExchangeRatesFromCB(country string) (models.CachedRates, error)
	ConvertRate(country string, from string, to string, value string) (string, error)
}
