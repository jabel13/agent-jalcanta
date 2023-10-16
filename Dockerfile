# Start the Go app build
FROM golang:latest AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Download dependencies
RUN go mod download

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Print working directory
RUN pwd && find .


CMD ["./main", "-poll=120"]

