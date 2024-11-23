FROM golang:1.21

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o edi_gateway ./cmd/main.go

CMD ["./edi_gateway"]