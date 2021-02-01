FROM golang:1.15-alpine AS builder

WORKDIR /sso

RUN apk add --no-cache git bash

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Verify dependencies
RUN go mod verify

COPY . .

RUN go build -o build/sso .




FROM alpine:3.13.0 AS production

RUN apk add --no-cache bash

WORKDIR /sso

COPY --from=builder /sso/migration/*.sql ./migration/
COPY --from=builder /sso/build/sso ./build/sso
EXPOSE 8080
ENTRYPOINT /sso/build/sso start

