package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Outcome struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Market struct {
	Key      string    `json:"key"`
	Outcomes []Outcome `json:"outcomes"`
}

type Bookmaker struct {
	Title   string  `json:"title"`
	Markets []Market `json:"markets"`
}

func main() {

	apiKey := "e74a90247906a097ffa99c9a4a611344"
	apiUrl := "https://api.the-odds-api.com/v4/sports/baseball_mlb/odds/?apiKey=" + apiKey

	response, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Check if the response status code is not 200 OK
	if response.StatusCode != http.StatusOK {
		fmt.Printf("API returned a non-OK status code: %d\n", response.StatusCode)
		return
	}

	// Decode the JSON response into a slice of Bookmaker structs
	var bookmakers []Bookmaker
	err = json.NewDecoder(response.Body).Decode(&bookmakers)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Loop through the games and print the bookmaker title, markets key, and outcomes
	for _, bookmaker := range bookmakers {
		fmt.Println("Bookmaker Title:", bookmaker.Title)

		for _, market := range bookmaker.Markets {
			fmt.Println("Market Key:", market.Key)

			fmt.Println("Outcomes:")
			for _, outcome := range market.Outcomes {
				fmt.Printf("  Name: %s, Price: %d\n", outcome.Name, outcome.Price)
			}
		}
		fmt.Println() // Add an empty line between each game
	}
}