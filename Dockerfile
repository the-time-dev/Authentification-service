FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/auth-service .
COPY . .

EXPOSE 8080

CMD ["./auth-service"] 