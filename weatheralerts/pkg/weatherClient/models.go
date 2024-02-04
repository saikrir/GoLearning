package weatherClient

import (
	"time"
)

type WeatherEntry struct {
	DateTxt string `json:"dt_txt"`
	Date    int64  `json:"dt"`
	Main    struct {
		Temp       float32 `json:"temp"`
		FeelsLike  float32 `json:"feels_like"`
		MinTemp    float32 `json:"temp_min"`
		MaxTemp    float32 `json:"temp_max"`
		Pressure   int32   `json:"pressure"`
		SeaLevel   int32   `json:"sea_level"`
		GrondLevel int32   `json:"grnd_level"`
		Humidity   int32   `json:"humidity"`
	} `json:"main"`

	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`

	Clouds struct {
		All int32 `json:"all"`
	} `json:"clouds"`

	Wind struct {
		Speed float32 `json:"speed"`
		Deg   int32   `json:"deg"`
		Gust  float32 `json:"gust"`
	} `json:"wind"`
}

type WeatherResponse struct {
	WeatherEntries []WeatherEntry `json:"list"`
	City           struct {
		Name       string `json:"name"`
		Country    string `json:"country"`
		Sunset     int64  `json:"sunset"`
		Sunrise    int64  `json:"sunrise"`
		Population int64  `json:"population"`
	} `json:"city"`
}

func (w *WeatherResponse) GetSunrise() time.Time {
	return time.Unix(w.City.Sunrise, 0)
}

func (w *WeatherResponse) GetSunset() time.Time {
	return time.Unix(w.City.Sunset, 0)
}

func (w *WeatherResponse) GroupByDay() map[string][]WeatherEntry {
	groupedDays := make(map[string][]WeatherEntry)

	for _, wEntry := range w.WeatherEntries {
		timeEntry := time.Unix(wEntry.Date, 0).Format("2006-02-01")
		entries, found := groupedDays[timeEntry]
		if !found {
			entries = []WeatherEntry{wEntry}
		} else {
			entries = append(entries, wEntry)
		}

		groupedDays[timeEntry] = entries
	}

	return groupedDays
}
