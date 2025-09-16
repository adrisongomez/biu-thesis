FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./custom-collector ./cmd/collectors

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/custom-collector /custom-collector
EXPOSE 4317
CMD ["/custom-collector"]