FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/api ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/bin/api ./api
COPY openapi.yaml ./openapi.yaml

EXPOSE 8080

CMD ["/app/api"]
