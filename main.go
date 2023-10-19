package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	loggly "github.com/jamespearly/loggly"
	"flag"
	"os"
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

func getAPIResponse(apiUrl string) (*http.Response, error) {
    response, err := http.Get(apiUrl)
    if err != nil {
        return nil, err
    }

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned a non-OK status code: %d", response.StatusCode)
    }
    return response, nil
}

func readAndParseResponse(response *http.Response) ([]Game, int, error) {
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, 0, err
    }
    
	contentSize := len(body)

    var games []Game
    err = json.Unmarshal(body, &games)

    if err != nil {
        return nil, 0, err
    }
    return games, contentSize, nil
}

func printGameDetails(games []Game) {
    for _, game := range games {
        fmt.Println("Game ID:", game.ID)
        for _, bookmaker := range game.Bookmakers {
            fmt.Println()
            fmt.Println("Bookmaker Title:", bookmaker.Title)
            for _, market := range bookmaker.Markets {
                fmt.Println("Outcomes:")
                for _, outcome := range market.Outcomes {
                    fmt.Printf("  Name: %s, Price: %d\n", outcome.Name, outcome.Price)
                }
            }
        }
        fmt.Println() // Add an empty line between each game
    }
}


func proccessMLBOdds() {

	var tag string
	tag = "Sports-Betting-Agent"

	client := loggly.New(tag)

	apiKey := "e74a90247906a097ffa99c9a4a611344"
	apiUrl := "https://api.the-odds-api.com/v4/sports/basketball_nba/odds/?apiKey=" + apiKey + "&regions=us" + "&markets=h2h" + "&oddsFormat=american"

    // Make the HTTP GET request with the updated URL
    response, err := getAPIResponse(apiUrl)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer response.Body.Close()


	games, contentSize, err := readAndParseResponse(response)

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	printGameDetails(games)

    // Use Sprintf to format string 
    formattedMsg := fmt.Sprintf("Size of JSON content: %d bytes", contentSize)
    err = client.EchoSend("info", formattedMsg)
    fmt.Println("err:", err)

}

func main() {

	// Define a new integer flag polling interval with default value 120
	// The user can specify the polling interval with -poll=<minutes>
	poll := flag.Int("poll", 0, "Polling interval in minutes")

	// Parse the flag
	flag.Parse()

    // Check if 'arg' is provided
	if *poll == 0 {
		fmt.Println("Error: The '-poll' flag is required.")
		os.Exit(1) // Exit with a non-zero status code to indicate an error
	}

	for {
		proccessMLBOdds()

		// Sleep for the specified duration
		time.Sleep(time.Duration(*poll) * time.Minute)
	}

}