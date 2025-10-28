FROM golang:1.23.9-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o az-dump main.go

FROM mcr.microsoft.com/azure-cli:2.78.0
LABEL org.opencontainers.image.source=https://github.com/nu12/az-dump

WORKDIR /app

COPY --from=builder /app/az-dump /app/az-dump

ENTRYPOINT ["./az-dump"]