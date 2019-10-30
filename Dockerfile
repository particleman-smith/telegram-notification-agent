# Builder stage
FROM golang as builder
VOLUME /var/log/telegram-notification-agent
COPY backend /go/src/github.com/particleman-smith/telegram-notification-agent
RUN go get -v -u github.com/gorilla/mux && go get -v -u github.com/lib/pq && go get -v -u github.com/rs/cors && go get -v -u github.com/go-telegram-bot-api/telegram-bot-api
RUN go install telegram-notification-agent

# Executable stage
FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/telegram-notification-agent /app/telegram-notification-agent
EXPOSE 9090
ENTRYPOINT ./telegram-notification-agent
