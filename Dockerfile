FROM golang:1.26-alpine

WORKDIR /app

# Install Air for development auto reload.
RUN go install github.com/air-verse/air@latest

# Make sure Go-installed binaries are available.
ENV PATH="/go/bin:${PATH}"

COPY go.mod ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "."]