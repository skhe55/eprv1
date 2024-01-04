package models

import (
	"encoding/xml"
	"time"
)

type Currency struct {
	NumCode   int    `xml:"NumCode"`
	Nominal   int    `xml:"Nominal"`
	CharCode  string `xml:"CharCode"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

type Valute []struct {
	XMLName xml.Name `xml:"Valute"`
	Currency
}

type XMLDocument struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date"`
	Valute  Valute   `xml:"Valute"`
}

type CachedRates struct {
	Rates      []Rates
	MapOfRates map[string]string
	LastUpdate time.Time
}

type Rates struct {
	CharCode string
	Value    string
}
