FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /var/back

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o main ./balancer/cmd/main.go

FROM alpine:edge as prod

RUN apk add --no-cache bash

WORKDIR /app

COPY --from=build /var/back/main .
COPY --from=build /var/back/.env .
COPY --from=build /var/back/backend/configs ./configs

EXPOSE 8080
ENTRYPOINT ["./main"]