FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/server /server
COPY --from=builder /app/cmd/server/.env.custom-collector /.env
EXPOSE 5000
CMD ["/server"]