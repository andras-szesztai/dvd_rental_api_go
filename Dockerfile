FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o ./bin/main ./cmd/api
  
# Install migrate CLI in your Dockerfile
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.tar.gz -C /usr/local/bin && \
    chmod +x /usr/local/bin/migrate && \
    rm migrate.tar.gz

# Run migrations using environment variable
CMD migrate -path /app/migrations -database "$DB_ADDR" up && ./bin/main