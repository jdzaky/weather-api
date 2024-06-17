package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rs/cors"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:OpenWeatherMapApiKey`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiKey")
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city + "&units=metric")
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	// city := strings.SplitN(r.URL.Path, "/", 3)[2]
	city := r.URL.Query().Get("city")
	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", handleWeather)

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	handler := cors.Default().Handler(http.DefaultServeMux)

	log.Println("Server started on :9999")
	log.Fatal(http.ListenAndServe(":9999", handler))

	http.ListenAndServe(":9999", nil)
}
