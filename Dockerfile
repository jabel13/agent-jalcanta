# Start the Go app build
FROM golang:latest AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Download dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

CMD ["./main"]

