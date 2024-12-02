# Start from the official Go image
FROM golang:1.22

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Generate Swagger docs
RUN swag init

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the executable
CMD ["./main"]