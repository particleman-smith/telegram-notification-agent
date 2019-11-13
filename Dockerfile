# Builder stage
FROM golang as builder
COPY app /go/src/github.com/the-codesmith/telegram-notification-agent/app
WORKDIR /go/src/github.com/the-codesmith/telegram-notification-agent/app
RUN go get -v -u github.com/gorilla/mux && go get -v -u github.com/lib/pq && go get -v -u github.com/rs/cors && go get -v -u github.com/go-telegram-bot-api/telegram-bot-api
RUN go install -tags netgo -a -v

# Executable stage
FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /app/
COPY --from=builder /go/bin/app /app/telegram-notification-agent
EXPOSE 9090
ENTRYPOINT ./telegram-notification-agent
