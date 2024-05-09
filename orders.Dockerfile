FROM golang:1.22.3

WORKDIR /order_service
COPY . .
RUN go mod tidy
RUN go build -o orders examples/services/orders/cmd/main.go

EXPOSE 5002
CMD ["./orders"]