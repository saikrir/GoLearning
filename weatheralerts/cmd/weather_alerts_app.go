package cmd

import (
	"fmt"
	"os"
	"weatheralerts/pkg/appLogger"
	"weatheralerts/pkg/weatherClient"
)

func getApiKey() string {
	apiKey := os.Getenv("API_KEY")

	if len(apiKey) == 0 {
		appLogger.Fatal("API key was not found")
		return ""
	}

	return apiKey
}

func LaunchApp() {
	appLogger.Info("AppLaunch", "ApiKey", getApiKey())
	var lat, lon float64 = 29.5612, -98.6802

	weatherResponse, err := weatherClient.GetWeatherForNextFiveDays(lat, lon, getApiKey())
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	fmt.Println(weatherResponse.GroupByDay()["2024-04-02"])
}
