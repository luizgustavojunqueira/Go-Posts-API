FROM golang:1.23.0

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/main ./cmd/main.go


CMD ["./bin/main"]
