package weatherClient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"weatheralerts/pkg/appLogger"
)

type WeatherDataProvider interface {
	getWeatherData() (io.ReadCloser, error)
}

type ApiWeatherDataProvider struct {
	url string
}

func apiWeatherDataSource(lat, lon float64, apiKey string) WeatherDataProvider {

	const apiUrl = "http://api.openweathermap.org/data/2.5/forecast"

	url, err := url.Parse(apiUrl)
	if err != nil {
		appLogger.Fatal("Failed to prase url " + apiUrl)
		return nil
	}

	values := url.Query()
	values.Add("lat", strconv.FormatFloat(lat, 'f', 4, 64))
	values.Add("lon", strconv.FormatFloat(lon, 'f', 4, 64))
	values.Add("appid", apiKey)
	values.Add("units", "imperial")

	url.RawQuery = values.Encode()

	return &ApiWeatherDataProvider{url: url.String()}
}

type FileWeatherDataProvider struct {
	filePath string
}

func fileWeatherDataSouce(filePath string) WeatherDataProvider {
	return &FileWeatherDataProvider{filePath: filePath}
}

func (wp ApiWeatherDataProvider) getWeatherData() (io.ReadCloser, error) {

	resp, err := http.Get(wp.url)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (fp FileWeatherDataProvider) getWeatherData() (io.ReadCloser, error) {
	dataFile, err := os.OpenFile(fp.filePath, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	return dataFile, nil
}

func GetWeatherForNextFiveDays(lat float64, lon float64, apiKey string) (*WeatherResponse, error) {
	//url := getWeatherDataUrl(lat, lon, apiKey)

	//weatherProvider := ApiWeatherDataProvider{url: url.String()}

	const filePath = "./Test.json"

	weatherProvider := fileWeatherDataSouce(filePath)

	data, err := weatherProvider.getWeatherData()

	if err != nil {
		appLogger.Error("Failed to call Weather API ", err.Error())
		return nil, err
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(data).Decode(&weatherResponse)

	if err != nil {
		appLogger.Error("Failed to unmarshall JSON", err.Error())
		return nil, err
	}

	defer data.Close()

	return &weatherResponse, nil
}
