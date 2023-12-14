FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./internal ./internal
COPY ./cmd/main.go ./cmd/main.go

# Build the Go app
RUN go build -o app ./cmd/main.go

FROM alpine as main

COPY --from=builder /app/app /bin/app
COPY ./config ./config

ENV CONFIG_PATH=/config/application-config.yaml

ENTRYPOINT [ "/bin/app" ]
