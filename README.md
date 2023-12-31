# Sports Odds Data Tracker in Go

## Description
This application is developed in Go and is designed to fetch, parse, and store NBA betting odds data from [The-Odds-API](https://the-odds-api.com/). It integrates with AWS DynamoDB to store the data and is designed to periodically poll the API, ensuring that the data remains up-to-date. The application is Dockerized for straightforward deployment.

## Features
- Fetches NBA odds data from The-Odds-API using an API key.
- Parses the retrieved data and structures it into Go data structures.
- Stores the processed data in AWS DynamoDB.
- Logs the process and any errors to Loggly for monitoring.
- Runs in Docker's virtual containers. 
- Customizable polling interval for data fetching.

## Prerequisites
- Go programming environment (version 1.x or later)
- Docker installed on your machine
- AWS account with DynamoDB access
- The-Odds-API key (emailed to you when you sign up for a plan)
- Loggly account for logging

## Setup
1. **Clone the Repository**: Clone the repository to your local environment using `git clone <repository_url>`
2. **Navigate to the Project Directory**: Change directory to the project folder with `cd <repository_name>`
3. **Install Dependencies**: Run `go mod tidy` to install the necessary Go dependencies.
4. **Build Docker Image**: Build the Docker image using the command `sudo docker build -t agent-jalcanta .`
5. **Run Docker Container**: Set your environment variables and start the Docker container with the command: <br>`docker run -e AWS_ACCESS_KEY_ID=<Your_Access_Key_ID> -e AWS_SECRET_ACCESS_KEY=<Your_Secret_Access_Key> -e AWS_DEFAULT_REGION=<Your_AWS_Region> -e API_KEY=<Your_API_Key> -e LOGGLY_TOKEN=<Your_Loggly_Token> agent-jalcanta` <br>replacing the placeholders with your actual AWS credentials, API key, and Loggly token.

## Usage
After starting the Docker container, the application will begin fetching and processing NBA betting odds data. By default, it polls the API every 2 hours. However, you can customize this interval using the -poll flag when running the Docker container. To adjust the polling interval, append -poll to your Docker run command. For instance, for a 10-minute interval, use<br>`docker run -e AWS_ACCESS_KEY_ID=<Your_Access_Key_ID> -e AWS_SECRET_ACCESS_KEY=<Your_Secret_Access_Key> -e AWS_DEFAULT_REGION=<Your_AWS_Region> -e API_KEY=<Your_API_Key> -e LOGGLY_TOKEN=<Your_Loggly_Token> agent-jalcanta -poll=10` <br>The console and Loggly dashboard can be used to monitor the application's operations. To stop the application, use the appropriate Docker commands for stopping or pausing the container.
