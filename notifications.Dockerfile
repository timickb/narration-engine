FROM golang:1.22.3

WORKDIR /notification_service
COPY . .
RUN go mod tidy
RUN go build -o notifications examples/services/notifications/cmd/main.go

EXPOSE 5001
CMD ["./notifications"]