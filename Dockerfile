# Start with the Go base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Install dockerize 
RUN wget https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz && \
    tar -C /usr/local/bin -xzvf dockerize-linux-amd64-v0.6.1.tar.gz && \
    rm dockerize-linux-amd64-v0.6.1.tar.gz

# Build the Go application
RUN go build -o main ./cmd/server/main.go

# Expose the application port
EXPOSE 8080


# Command to run the application, waiting for MySQL and Elasticsearch
CMD ["dockerize", "-wait", "tcp://db:3306", "-timeout", "30s", "-wait", "tcp://elasticsearch:9200", "-timeout", "30s", "./main","--host=0.0.0.0", "--port=8080"]