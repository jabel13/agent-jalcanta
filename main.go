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
    "strconv"
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

// Sends an HTTP GET request to the specified apiUrl
// and returns a pointer to the HTTP response and any error encountered
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

func readAndParseResponse(response *http.Response) ([]Game, error) {

    // Initialize loggly tag and client
    var tag string
	tag = "Sports-Betting-Agent"

	client := loggly.New(tag)

    // reads the entire content of the HTTP response body into 'body' variable
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }
    
    contentSize := len(body)

    // Attempt to parse the 'body' (response in JSON) into a slice of Game objects
    var allGames []Game
    // &allGames is a pointer to the 'allGames' variable 
    // this is where the parsed data will be stored after decoding from JSON
    err = json.Unmarshal(body, &allGames)
    if err != nil {
        return nil, err
    }

    // Limit the number of Game objects to a maximum of 7 games
    maxGames := 7
    if len(allGames) > maxGames {
        allGames = allGames[:maxGames]
    }

    // Use Sprintf to format string 
    formattedMsg := fmt.Sprintf("Size of JSON content: %d bytes", contentSize)
    err = client.EchoSend("info", formattedMsg)
    fmt.Println("err:", err)

    return allGames, nil
}

// Takes a slice of Game objects - iterate through each game and print details
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
// Must specify AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_DEFAULT_REGION 
// as environmnet variables 
// Create a new AWS session. If the session creation fails,
// the program will log the error and terminate.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
    })
    if err != nil {
        log.Fatalf("Failed to create AWS session: %s", err)
    }

    
    // Create DynamoDB client
    svc := dynamodb.New(sess)


    // Create table
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

                // Marshal the DynamoItem into a map so it can be saved to DynamoDB
				item, err := dynamodbattribute.MarshalMap(dynamoItem)
				if err != nil {
					log.Fatalf("Got error marshalling map: %s", err)
				}
                
                // Create the PutItem input structure based on the marshaled item
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

	// Fetch the API key from environment variables
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalf("Error: API_KEY environment variable not set")
		return
	}
	apiUrl := "https://api.the-odds-api.com/v4/sports/basketball_nba/odds/?apiKey=" + apiKey + "&regions=us" + "&markets=h2h" + "&oddsFormat=american"

    // Make the HTTP GET request with the updated URL
    response, err := getAPIResponse(apiUrl)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }


	games, err := readAndParseResponse(response)

    if err != nil {
        fmt.Println("Error:", err)
        response.Body.Close()
        return
    }

    response.Body.Close()

	printGameDetails(games)

    // Write data to DynamoDB
    err = writeToDynamoDB(games)
    if err != nil {
        fmt.Println("Error writing to DynamoDB:", err)
        return
    }

}

func main() {

	// Retrieve the environment variable
	pollStr := os.Getenv("POLL_INTERVAL")

    // Convert the environment variable to an integer
	poll, err := strconv.Atoi(pollStr)
	if err != nil {
		fmt.Println("Error: Invalid POLL_INTERVAL value. It should be an integer.")
		os.Exit(1)
	}

	// Check if poll is provided and is a positive value
	if poll <= 0 {
		fmt.Println("Error: POLL_INTERVAL must be a positive integer.")
		os.Exit(1)
	}


	for {
		proccessNbaOdds()

		// Sleep for the specified duration
		time.Sleep(time.Duration(poll) * time.Minute)
	}

}