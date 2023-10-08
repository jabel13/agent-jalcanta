package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	loggly "github.com/jamespearly/loggly"
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
    Key     string   `json:"key"`
    Title   string   `json:"title"`
    Markets []Market `json:"markets"`
}

type Game struct {
    ID            string      `json:"id"`
    Bookmakers    []Bookmaker `json:"bookmakers"`
}

func main() {

	var tag string
	tag = "Sports-Betting-Server"

	client := loggly.New(tag)

	apiKey := "e74a90247906a097ffa99c9a4a611344"
	apiUrl := "https://api.the-odds-api.com/v4/sports/baseball_mlb/odds/?apiKey=" + apiKey + "&regions=us" + "&markets=h2h" + "&oddsFormat=american"

    // Make the HTTP GET request with the updated URL
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

	// Read the response body as a string
// body, err := ioutil.ReadAll(response.Body)
// if err != nil {
//     fmt.Println("Error reading response body:", err)
//     return
// }

// Print the raw JSON response
// fmt.Println("Raw JSON Response:", string(body))

	// Decode the JSON response into a slice of Bookmaker structs
	var games []Game
	err = json.NewDecoder(response.Body).Decode(&games)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Valid EchoSend (message echoed to console and no error returned)
	err = client.EchoSend("debug", "This is a debug message")
	fmt.Println("err:", err)

	// Valid Send (no error returned)
	err = client.Send("info", "Message received")
	fmt.Println("err:", err)

	fmt.Println(response.ContentLength)

	// for _, game := range games {
	// 	fmt.Println("Game ID:", game.ID)
	// 	// Loop through the games and print the bookmaker title, markets key, and outcomes
	// 	for _, bookmaker := range game.Bookmakers {
	// 		fmt.Println()
	// 		fmt.Println("Bookmaker Title:", bookmaker.Title)

	// 		for _, market := range bookmaker.Markets {

	// 			fmt.Println("Outcomes:")
	// 			for _, outcome := range market.Outcomes {
	// 				fmt.Printf("  Name: %s, Price: %d\n", outcome.Name, outcome.Price)
	// 			}
	// 		}
	// 	}
	// 	fmt.Println() // Add an empty line between each game
	// }
}