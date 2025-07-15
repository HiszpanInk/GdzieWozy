package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// struct for configuration
type Config struct {
	API_KEY string
}

type apiDependent struct {
	api_key string
}

// trzeba tu zrobić żeby mieć zapytanie na konkretne linie, trza też wykombinować żeby można
// accessować tylko z poziomu serwera
// todo: cache'ing
func (apiDependent apiDependent) buses(w http.ResponseWriter, req *http.Request) {

	answer, err := http.Get(fmt.Sprintf("https://api.um.warszawa.pl/api/action/busestrams_get/?resource_id=f2e5503e927d-4ad3-9500-4ab9e55deb59&apikey=%s&type=1", apiDependent.api_key))
	req.ParseForm()
	//lines := req.Form.Get("lines")
	log.Print(url.Values{}["lines"])
	data, err := io.ReadAll(answer.Body)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("cos jest")
		fmt.Fprintf(w, string(data))
	}
}

func wawmap(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "main.html")
}

func main() {
	//get API key to https://api.um.warszawa.pl
	filepath := "./config.json"
	data, err := os.ReadFile(filepath)
	var config Config

	if err != nil {
		fmt.Println("Error reading config file", err)
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error processing config file", err)
		return
	}
	apiDependent := apiDependent{api_key: config.API_KEY}

	http.HandleFunc("/buses", apiDependent.buses)
	http.HandleFunc("/map", wawmap)

	http.ListenAndServe(":8090", nil)
}
