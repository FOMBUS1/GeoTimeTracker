package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiKey string
	apiURL string
}

func NewClient(apiKey, apiURL string) *Client {
	return &Client{apiKey: apiKey, apiURL: apiURL}
}

func (c *Client) Fetch(ctx context.Context, lat, long float32) (string, error) {
	address := "Адрес не найден"
	resp, err := http.Get(c.apiURL)
	if err != nil {
		return address, fmt.Errorf("ошибка запроса к Яндекс.Картам: %v", err)
	}
	defer resp.Body.Close()

	var geoData YandexGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return address, fmt.Errorf("ошибка парсинга ответа Яндекса: %v", err)
	}

	if len(geoData.Response.GeoObjectCollection.FeatureMember) > 0 {
		address = geoData.Response.GeoObjectCollection.FeatureMember[0].GeoObject.MetaDataProperty.GeocoderMetaData.Address.Formatted
	}

	return address, nil
}
