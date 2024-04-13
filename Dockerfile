FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM alpine:latest
RUN apk update && apk add --no-cache postgresql-client bash
COPY --from=builder /app/app /app/app
COPY ./config/config.yaml /app/config/config.yaml
COPY ./wait-for-it.sh /app/wait-for-it.sh
COPY ./wait-for-it.sh /wait-for-it.sh
WORKDIR /app
ENV CONFIG_PATH=/app/config/config.yaml
CMD ["./app"]
