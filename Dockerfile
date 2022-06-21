FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./

RUN go build -o bin/ipinfo cmd/api/main.go

FROM alpine

WORKDIR /

COPY --from=build /app/bin/ipinfo /ipinfo

EXPOSE 4000

ENTRYPOINT ["/ipinfo"]