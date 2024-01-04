package usecase

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"erpv1/internal/converter"
	"erpv1/internal/models"
	myErrors "erpv1/internal/server/errors"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

type converterUseCase struct {
	converterRepo converter.Repository
}

func NewConverterUseCase(converterRepo converter.Repository) converter.UseCase {
	return &converterUseCase{converterRepo: converterRepo}
}

func (u *converterUseCase) UpdateExchangeRatesFromRUCB() (models.CachedRates, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://www.cbr.ru/scripts/XML_daily.asp", nil)
	if err != nil {
		fmt.Printf("Failed to create new request, error: %v", err)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 OPR/105.0.0.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed get data from cb, error: %v", err)
		return models.CachedRates{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var root models.XMLDocument
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&root)
	if err != nil {
		fmt.Printf("Failed decode data from cb, error: %v", err)
		return models.CachedRates{}, err
	}

	return u.converterRepo.SetExchangeRatesFromRUCB(root), nil
}

func (u *converterUseCase) UpdateExchangeRatesFromThaiCB() (models.CachedRates, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://apigw1.bot.or.th/bot/public/Stat-ExchangeRate/v2/DAILY_AVG_EXG_RATE/?start_period=2024-01-03&end_period=2024-01-03", nil)
	if err != nil {
		fmt.Printf("Failed to create new request, error: %v", err)
		return models.CachedRates{}, err
	}

	req.Header.Set("X-IBM-Client-Id", "c2bbe063-d0ff-456c-bc08-fbd5115fb340")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed get data from cb, error: %v", err)
		return models.CachedRates{}, err
	}
	var data models.ThaiCBResponse
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return u.converterRepo.SetExchangeRatesFromThaiCB(data), nil
}

func (u *converterUseCase) GetExchangeRatesFromCB(country string) (models.CachedRates, error) {
	var callbackForUpdateRates func() (models.CachedRates, error)
	var reqInterval time.Duration

	upperCountry := strings.ToUpper(country)

	if upperCountry == "RU" {
		callbackForUpdateRates = u.UpdateExchangeRatesFromRUCB
		reqInterval = time.Minute
	} else if upperCountry == "THAI" {
		callbackForUpdateRates = u.UpdateExchangeRatesFromThaiCB
		reqInterval = time.Hour
	}

	data, err := u.converterRepo.GetExchangeRatesFromCB(country, reqInterval)

	if errors.Is(err, myErrors.ErrNeededUpdate) {
		data, err := callbackForUpdateRates()
		if err != nil {
			return models.CachedRates{}, err
		}
		return data, nil
	}
	if errors.Is(err, myErrors.ErrNotFound) {
		data, err := callbackForUpdateRates()
		if err != nil {
			return models.CachedRates{}, err
		}
		return data, nil
	}
	return data, nil
}

func (u *converterUseCase) ConvertRate(country string, from string, to string, value string) (string, error) {
	upperCountry, upperTo, upperFrom := strings.ToUpper(country), strings.ToUpper(to), strings.ToUpper(from)
	if upperTo == upperFrom {
		return "", errors.New("params from and to must be different")
	}
	if upperCountry == "RU" {
		data, err := u.GetExchangeRatesFromCB(country)
		if err != nil {
			return "", err
		}
		if upperTo == "RUB" {
			valuteValue, ok := data.MapOfRates[upperFrom]
			if !ok {
				return "", myErrors.ErrInvalidCurrency
			}
			rate, _ := strconv.ParseFloat(strings.Replace(valuteValue, ",", ".", -1), 32)
			v, _ := strconv.ParseInt(value, 10, 32)
			return strconv.FormatFloat(float64(v)*float64(rate), 'f', -1, 64), nil
		} else {
			if upperFrom != "RUB" {
				nativeCurrencyValue, ok := data.MapOfRates[upperFrom]
				if !ok {
					return "", myErrors.ErrInvalidCurrency
				}
				destCurrencyValue, ok := data.MapOfRates[upperTo]
				if !ok {
					return "", myErrors.ErrInvalidCurrency
				}
				nativeRate, _ := strconv.ParseFloat(strings.Replace(nativeCurrencyValue, ",", ".", -1), 32)
				destRate, _ := strconv.ParseFloat(strings.Replace(destCurrencyValue, ",", ".", -1), 32)
				v, _ := strconv.ParseInt(value, 10, 32)

				amountInNativeCurrency := float64(v) * float64(nativeRate)

				return strconv.FormatFloat(amountInNativeCurrency/destRate, 'f', -1, 64), nil
			} else {
				valuteValue, ok := data.MapOfRates[upperTo]
				if !ok {
					return "", myErrors.ErrInvalidCurrency
				}
				rate, _ := strconv.ParseFloat(strings.Replace(valuteValue, ",", ".", -1), 32)
				v, _ := strconv.ParseInt(value, 10, 32)
				return strconv.FormatFloat(float64(v)/rate, 'f', -1, 64), nil
			}
		}
	} else if upperCountry == "THAI" {
		data, err := u.GetExchangeRatesFromCB(country)
		if err != nil {
			return "", err
		}
		if upperTo == "THB" {
			valuteValue, ok := data.MapOfRates[upperFrom]
			if !ok {
				return "", myErrors.ErrInvalidCurrency
			}
			rate, _ := strconv.ParseFloat(strings.Replace(valuteValue, ",", ".", -1), 32)
			v, _ := strconv.ParseInt(value, 10, 32)
			return strconv.FormatFloat(float64(v)*rate, 'f', -1, 64), nil
		} else {
			if upperFrom != "THB" {
				nativeCurrencyValue, ok := data.MapOfRates[upperFrom]
				if !ok {
					return "", myErrors.ErrInvalidCurrency
				}
				destCurrencyValue, ok := data.MapOfRates[upperTo]
				if !ok {
					return "", myErrors.ErrInvalidCurrency
				}
				nativeRate, _ := strconv.ParseFloat(strings.Replace(nativeCurrencyValue, ",", ".", -1), 32)
				destRate, _ := strconv.ParseFloat(strings.Replace(destCurrencyValue, ",", ".", -1), 32)
				v, _ := strconv.ParseInt(value, 10, 32)

				amountInNativeCurrency := float64(v) * float64(nativeRate)

				return strconv.FormatFloat(amountInNativeCurrency/destRate, 'f', -1, 64), nil
			} else {
				destCurrencyValue, ok := data.MapOfRates[upperTo]
				if !ok {
					return "", myErrors.ErrInvalidCurrency
				}
				destRate, _ := strconv.ParseFloat(strings.Replace(destCurrencyValue, ",", ".", -1), 32)
				v, _ := strconv.ParseInt(value, 10, 32)

				return strconv.FormatFloat(float64(v)/destRate, 'f', -1, 64), nil
			}
		}
	}

	return "", nil
}
