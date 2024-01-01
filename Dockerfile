FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go build -v -o ./server ./cmd/server/

FROM alpine
WORKDIR /app
COPY .env.prod /app/.env
COPY --from=builder /app/server ./server
CMD ./server