package cryptocompare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	path       = "https://min-api.cryptocompare.com/data"
	allCryptos = "pricemulti"
	fsyms      = "fsyms"
	tsyms      = "tsyms"
	argsSep    = ","
	pathSep    = "/"
	costIn     = "USD"
)

type CryptoCompare struct {
}

func (c *CryptoCompare) GetAll(titles []string, in string) (map[string]float64, error) {
	var resultsRaw = map[string]map[string]interface{}{}

	rawURL, err := url.Parse(strings.Join([]string{path, allCryptos}, pathSep))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add(fsyms, strings.Join(titles, argsSep))
	params.Add(tsyms, strings.Join([]string{in}, argsSep))
	rawURL.RawQuery = params.Encode()

	res, err := http.Get(rawURL.String())
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &resultsRaw); err != nil {
		return nil, err
	}

	return c.castResultData(resultsRaw)
}

func (c *CryptoCompare) GetSpecial(title string, in string) (map[string]float64, error) {
	var resultsRaw = map[string]map[string]interface{}{}

	rawURL, err := url.Parse(strings.Join([]string{path, allCryptos}, pathSep))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add(fsyms, strings.Join([]string{title}, argsSep))
	params.Add(tsyms, strings.Join([]string{in}, argsSep))
	rawURL.RawQuery = params.Encode()

	res, err := http.Get(rawURL.String())
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &resultsRaw); err != nil {
		return nil, err
	}

	return c.castResultData(resultsRaw)
}

func (c *CryptoCompare) castResultData(in map[string]map[string]interface{}) (map[string]float64, error) {
	res := make(map[string]float64)
	for title, costMap := range in {
		cost, ok := costMap[costIn].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to assert: %v to float64 format", costMap[costIn])
		}
		res[title] = cost
	}
	return res, nil
}
