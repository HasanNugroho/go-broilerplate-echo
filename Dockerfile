# Use the official Golang image as the build stage
FROM golang:1.22 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create the workspace inside the container
WORKDIR /app

# # Copy go.mod and go.sum for dependency management
# COPY go.mod go.sum ./

# Copy all project files into the container
COPY . .

# RUN go mod tidy

# RUN go mod verify

# Download dependencies
RUN go mod download

# Build the application binary
RUN go build -o main cmd/main.go

# Use a minimal Alpine image for the final stage
FROM alpine:latest  

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=builder /app/main .

# Run the application when the container starts
CMD ["./main"]
