package converter

import (
	"erpv1/internal/models"
	"time"
)

type Repository interface {
	SetExchangeRatesFromRUCB(data models.XMLDocument) models.CachedRates
	SetExchangeRatesFromThaiCB(data models.ThaiCBResponse) models.CachedRates
	GetExchangeRatesFromCB(country string, reqInterval time.Duration) (models.CachedRates, error)
}
