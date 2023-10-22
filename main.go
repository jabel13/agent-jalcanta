package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"encoding/json"
	"fmt"
    "log"
	"net/http"
	"io/ioutil"
	"time"
	loggly "github.com/jamespearly/loggly"
	"flag"
	"os"
)

type Game struct {
	ID         string       `json:"id"`
	Bookmakers []Bookmaker  `json:"bookmakers"`
}

type Bookmaker struct {
	Key     string    `json:"key"`
    Title   string    `json:"title"`
	Markets []Market  `json:"markets"`
}

type Market struct {
	Outcomes []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type DynamoItem struct {
	GameID       string    `json:"id"`
	BookmakerKey string    `json:"key"`
	Outcomes     []Outcome `json:"outcomes"`
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

    var allGames []Game
    err = json.Unmarshal(body, &allGames)
    if err != nil {
        return nil, 0, err
    }

    // Limit to a maximum of 7 games
    maxGames := 7
    if len(allGames) > maxGames {
        allGames = allGames[:maxGames]
    }

    return allGames, contentSize, nil
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

func writeToDynamoDB(games []Game) error {
// Initialize a session that the SDK will use to load
// credentials from the shared credentials file ~/.aws/credentials
// and region from the shared configuration file ~/.aws/config.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
    })
    if err != nil {
        log.Fatalf("Failed to create AWS session: %s", err)
    }

    
    // Create DynamoDB client
    svc := dynamodb.New(sess)


    // Define the name of your table
    tableName := "nba-odds-jalcanta"

	// Iterate through each game and its bookmakers to create Dynamo items
	for _, game := range games {
		for _, bookmaker := range game.Bookmakers {
			for _, market := range bookmaker.Markets {
				dynamoItem := DynamoItem{
					GameID:       game.ID,
					BookmakerKey: bookmaker.Key,
					Outcomes:     market.Outcomes,
				}

				item, err := dynamodbattribute.MarshalMap(dynamoItem)
				if err != nil {
					log.Fatalf("Got error marshalling map: %s", err)
				}

				input := &dynamodb.PutItemInput{
					TableName: &tableName,
					Item:      item,
				}

				_, err = svc.PutItem(input)
				if err != nil {
					log.Fatalf("Got error calling PutItem: %s", err)
				}
			}
		}
	}

	return nil
}



func proccessNbaOdds() {

	var tag string
	tag = "Sports-Betting-Agent"

	client := loggly.New(tag)

	// Fetch the API key from environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("Error: API_KEY environment variable not set")
		return
	}
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

    // Write data to DynamoDB
    err = writeToDynamoDB(games)
    if err != nil {
        fmt.Println("Error writing to DynamoDB:", err)
        return
    }

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
		proccessNbaOdds()

		// Sleep for the specified duration
		time.Sleep(time.Duration(*poll) * time.Minute)
	}

}