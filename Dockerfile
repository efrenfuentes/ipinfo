FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -o bin/ipinfo cmd/api/main.go

EXPOSE 4000

CMD ["./bin/ipinfo"]