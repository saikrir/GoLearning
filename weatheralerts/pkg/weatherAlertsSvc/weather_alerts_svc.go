package weatherAlertsSvc

import "weatheralerts/pkg/weatherClient"

const HUMIDITY_LIMIT = 90.0
const WIND_SPEED_LIMIT = 18.0
const MAX_WEATHER_LIMIT = 100.0
const MIN_WEATHER_LIMIT = 40.0
const WEATHER_ID_RAIN = 500

type WeatherAlertsSvc struct {
	groupedMap map[string][]weatherClient.WeatherEntry
}

func NewWeatherAlertsSvc(weatherResponse *weatherClient.WeatherResponse) *WeatherAlertsSvc {
	weatherAlertsSvc := WeatherAlertsSvc{groupedMap: weatherResponse.GroupByDay()}
	return &weatherAlertsSvc
}

func (w *WeatherAlertsSvc) GatherWindAlerts() *string {
	for _, entries := range w.groupedMap {
		for _, entry := range entries {
			if entry.Wind.Speed >= WIND_SPEED_LIMIT {
				return &entry.DateTxt
			}
		}
	}
	return nil
}

func (w *WeatherAlertsSvc) GetRainAlert(weatherResponse weatherClient.WeatherResponse) *string {
	for _, entries := range w.groupedMap {
		for _, entry := range entries {
			for _, weather := range entry.WeatherPredications {
				if weather.Id == WEATHER_ID_RAIN {
					return &entry.DateTxt
				}
			}
		}
	}
	return nil
}
