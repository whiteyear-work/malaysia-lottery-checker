FROM golang:1.26-alpine

# Sets /app as the working folder inside the container.
WORKDIR /app

# Copies the Go module file first.
COPY go.mod ./

# download dependecy
RUN go mod download 

# Copies the entire project into the Docker image. from root and current
COPY . .

EXPOSE 8080

# Starts the Go application.
CMD ["go", "run", "."]