FROM golang:1.22.3

WORKDIR /payment_service
COPY . .
RUN go mod tidy
RUN go build -o payments cmd/main.go

EXPOSE 5002
CMD ["./payments"]