FROM golang:1.16-alpine as builder
WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o stewart-bot cmd/bot/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/stewart-bot .
COPY --from=builder /build/config.yaml .

ENV STEW_TOKEN=""
ENV STEW_SERVER_ID=""
ENV STEW_MAIN_CH=""
ENV STEW_VOICE_CH=""

ENTRYPOINT ["./stewart-bot"]