package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiKey string `json:"apiKey"`
	City   string `json:"city"`
	Units  string `json:"units"`
}

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func main() {
	const YELLOW = "\033[1;33m"
	const GREY = "\033[1;37m"
	const WHITE = "\033[1;37m"
	const BLUE = "\033[1;34m"
	const RESET = "\033[0m"

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executableDir := filepath.Dir(executablePath)
	envFilePath := filepath.Join(executableDir, "..", ".env")

	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	city := os.Getenv("CITY")
	units := os.Getenv("UNITS")

	if apiKey == "" || city == "" || units == "" {
		log.Fatal("Please set the API_KEY, CITY, and UNITS environment variables.")
	}

	var windSpeedUnits string
	if units == "metric" {
		windSpeedUnits = "m/s"
	} else if units == "imperial" {
		windSpeedUnits = "mph"
	} else {
		windSpeedUnits = "m/s"
	}

	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=%s&APPID=%s", city, units, apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		log.Fatal(err)
	}

	if weatherData.Cod != 200 {
		switch weatherData.Cod {
		case 404:
			fmt.Println("City not found")
			os.Exit(1)
		case 401:
			fmt.Println("Invalid API key")
			os.Exit(1)
		default:
			fmt.Println("ERROR")
			os.Exit(1)
		}
	}

	temp := strconv.FormatFloat(weatherData.Main.Temp, 'f', 0, 64)
	tempMin := strconv.FormatFloat(weatherData.Main.TempMin, 'f', 0, 64)
	tempMax := strconv.FormatFloat(weatherData.Main.TempMax, 'f', 0, 64)
	windSpeed := strconv.FormatFloat(weatherData.Wind.Speed, 'f', 0, 64)
	humidity := strconv.FormatFloat(float64(weatherData.Main.Humidity), 'f', 0, 64)

	println()

	switch weatherData.Weather[0].Main {
	case "Clear":
		fmt.Printf("     %s\\   /%s", YELLOW, RESET)
		fmt.Printf("      %s.-.%s\n", YELLOW, RESET)
		fmt.Printf("   %s‒ (   ) ‒%s", YELLOW, RESET)
		fmt.Printf("      %s`-᾿%s\n", YELLOW, RESET)
		fmt.Printf("     %s/   \\%s\n", YELLOW, RESET)
	case "Clouds":
		fmt.Printf("       %s.--.%s\n", GREY, RESET)
		fmt.Printf("    %s.-(    ).%s\n", GREY, RESET)
		fmt.Printf("   %s(___.__)__)%s\n", GREY, RESET)
	case "Rain":
		fmt.Printf("       %s.--.%s\n", BLUE, RESET)
		fmt.Printf("    %s.-(    ).%s\n", BLUE, RESET)
		fmt.Printf("   %s(___.__)__)%s\n", BLUE, RESET)
		fmt.Printf("    %sʻ‚ʻ‚ʻ‚ʻ‚ʻ%s\n", BLUE, RESET)
	case "Snow":
		fmt.Printf("       %s.--.%s\n", WHITE, RESET)
		fmt.Printf("    %s.-(    ).%s\n", WHITE, RESET)
		fmt.Printf("   %s(___.__)__)%s\n", WHITE, RESET)
		fmt.Printf("    %s* * * * *%s\n", WHITE, RESET)
	case "Thunderstorm":
		fmt.Printf("         .--.           Weather: storm\n")
		fmt.Printf("    %s.-(    ).%s       Temperature: %s\n", BLUE, RESET, temp)
		fmt.Printf("   %s(___.__)__)%s      Min/Max: %s/%s\n", BLUE, RESET, tempMin, tempMax)
	}

	fmt.Printf("\n")
	fmt.Printf("   Weather: clear\n")
	fmt.Printf("   Temperature: %s\n", temp)
	fmt.Printf("   Min/Max: %s/%s\n", tempMin, tempMax)
	fmt.Printf("   Wind speed: %s %s\n", windSpeed, windSpeedUnits)
	fmt.Printf("   Humidity: %s%%\n", humidity)

	println()
}
