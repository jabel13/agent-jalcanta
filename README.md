# NBA Betting Odds Processing Application in Go

## Description
This application is developed in Go and is designed to fetch, parse, and store NBA betting odds data from [The-Odds-API](https://the-odds-api.com/). It integrates with AWS DynamoDB to store the data and is designed to periodically poll the API, ensuring that the data remains up-to-date. The application is Dockerized for straightforward deployment and scalability.

## Features
- Fetches NBA odds data from The-Odds-API using an API token.
- Parses the retrieved data and structures it into Go data structures.
- Stores the processed data in AWS DynamoDB.
- Logs the process and any errors to Loggly for monitoring.
- Runs in a Docker container for easy deployment and scaling.
- Customizable polling interval for data fetching.

## Prerequisites
- Go programming environment (version 1.x or later)
- Docker installed on your machine
- AWS account with DynamoDB access
- The Odds API access token
- Loggly account for logging

## Setup
1. **Clone the Repository**: Clone the repository to your local environment using `git clone <repository_url>`.
2. **Navigate to the Project Directory**: Change directory to the project folder with `cd <repository_name>`.
3. **Install Dependencies**: If a `go.mod` file is present, run `go mod tidy` to install the necessary Go dependencies.
4. **Set Environment Variables**: Export your AWS credentials, The Odds API key, and Loggly token as environment variables.
5. **Build Docker Image**: Build the Docker image using the command `docker build -t agent-jalcanta .`.
6. **Run Docker Container**: Start the Docker container with the command: `docker run -e AWS_ACCESS_KEY_ID=<Your_Access_Key_ID> -e AWS_SECRET_ACCESS_KEY=<Your_Secret_Access_Key> -e AWS_DEFAULT_REGION=<Your_AWS_Region> -e API_KEY=<Your_API_Key> -e LOGGLY_TOKEN=<Your_Loggly_Token> -p 8080:8080 agent-jalcanta`, replacing the placeholders with your actual AWS credentials, API key, and Loggly token.

## Usage
After starting the Docker container, the application will automatically begin fetching and processing NBA betting odds data. By default, it polls the API every 2 hours. However, you can customize this interval using the -poll flag when running the Docker container. To adjust the polling interval, append -poll=<minutes> to your Docker run command, where <minutes> is the desired interval in minutes. For instance, for a 10-minute interval, use<br>`docker run -e AWS_ACCESS_KEY_ID=<Your_Access_Key_ID> -e AWS_SECRET_ACCESS_KEY=<Your_Secret_Access_Key> -e AWS_DEFAULT_REGION=<Your_AWS_Region> -e API_KEY=<Your_API_Key> -e LOGGLY_TOKEN=<Your_Loggly_Token> -p 8080:8080 agent-jalcanta -poll=10`. <br>The console and Loggly dashboard can be used to monitor the application's operations. To stop the application, use the appropriate Docker commands for stopping or pausing the container.
