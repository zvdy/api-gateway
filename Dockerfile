FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go get github.com/spf13/viper
RUN go build -o api-gateway ./cmd/main.go

EXPOSE 8080

CMD ["./api-gateway"]