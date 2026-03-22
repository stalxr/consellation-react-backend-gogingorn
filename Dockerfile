# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o charity-api ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/charity-api .
COPY --from=builder /app/openapi.yaml ./openapi.yaml

RUN mkdir -p uploads/images uploads/reports

EXPOSE 8080

CMD ["./charity-api"]
