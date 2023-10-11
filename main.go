package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
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
	apiUrl := "https://api.the-odds-api.com/v4/sports/baseball_mlb/odds/?apiKey=" + apiKey + "&regions=us" + "&markets=h2h" + "&oddsFormat=american"

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

	// Valid EchoSend (message echoed to console and no error returned)
	err = client.EchoSend("debug", "This is a debug message")
	fmt.Println("err:", err)

	// Valid Send (no error returned)
	err = client.Send("info", "Message received")
	fmt.Println("err:", err)

	fmt.Printf("Size of JSON content: %d bytes\n", contentSize)


}

func main() {

	for {
		proccessMLBOdds()

		// Sleep for 120 minutes before the next request
		// time.Sleep(120 * time.Minute)

		// Testing purposes
		time.Sleep(2 * time.Minute)
	}

}