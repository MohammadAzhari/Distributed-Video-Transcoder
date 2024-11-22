
FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /main main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /main /main
COPY app.env app.env

EXPOSE 8080

ENTRYPOINT ["/main"]