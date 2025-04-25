FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/auth-service .

EXPOSE 8080

CMD ["./auth-service"] 