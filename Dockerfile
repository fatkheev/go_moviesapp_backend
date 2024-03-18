FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go_moviesapp_backend ./cmd/main.go

EXPOSE 8081

CMD ["./go_moviesapp_backend"]