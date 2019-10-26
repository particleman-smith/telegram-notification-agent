# Builder stage
FROM golang as builder
VOLUME /var/log/telegram-notification-agent
COPY backend /go/src/telegram-notification-agent
RUN go get -v -u github.com/gorilla/mux && go get -v -u github.com/lib/pq && go get -v -u github.com/rs/cors
RUN go install scan-man

# Executable stage
FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/telegram-notification-agent /app/telegram-notification-agent
EXPOSE 9090
ENTRYPOINT ./scan-man