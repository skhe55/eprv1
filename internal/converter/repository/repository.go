package repository

import (
	"erpv1/internal/converter"
	"erpv1/internal/models"
	myErrors "erpv1/internal/server/errors"
	"erpv1/pkg/inmemdb"
	"time"
)

type converterInmemRepository struct {
	db *inmemdb.InmemDB
}

func NewConverterInmemRepository(db *inmemdb.InmemDB) converter.Repository {
	return &converterInmemRepository{db: db}
}

func (r *converterInmemRepository) SetExchangeRatesFromRUCB(data models.XMLDocument) models.CachedRates {
	var rates []models.Rates
	var mapOfRates = make(map[string]string)
	for _, valute := range data.Valute {
		rates = append(rates, models.Rates{
			CharCode: valute.CharCode,
			Value:    valute.Value,
		})
		mapOfRates[valute.CharCode] = valute.Value
	}

	cachedRates := models.CachedRates{
		Rates:      rates,
		MapOfRates: mapOfRates,
		LastUpdate: time.Now().UTC(),
	}

	r.db.Set("ru", cachedRates)

	return cachedRates
}

func (r *converterInmemRepository) SetExchangeRatesFromThaiCB(data models.ThaiCBResponse) models.CachedRates {
	var rates []models.Rates
	var mapOfRates = make(map[string]string)
	for _, valute := range data.Result.Data.DataDetail {
		rates = append(rates, models.Rates{
			CharCode: valute.CurrencyID,
			Value:    valute.Selling,
		})
		mapOfRates[valute.CurrencyID] = valute.Selling
	}

	cachedRates := models.CachedRates{
		Rates:      rates,
		MapOfRates: mapOfRates,
		LastUpdate: time.Now().UTC(),
	}

	r.db.Set("thai", cachedRates)

	return cachedRates
}

func (r *converterInmemRepository) GetExchangeRatesFromCB(country string, reqInterval time.Duration) (models.CachedRates, error) {
	value, ok := r.db.Get(country)
	if !ok {
		return models.CachedRates{}, myErrors.ErrNotFound
	}
	if time.Since(value.(models.CachedRates).LastUpdate) >= reqInterval {
		return models.CachedRates{}, myErrors.ErrNeededUpdate
	}
	return value.(models.CachedRates), nil
}
