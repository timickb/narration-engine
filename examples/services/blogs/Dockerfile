FROM golang:1.22.3

WORKDIR /blog_service
COPY . .
RUN go mod tidy
RUN go build -o blogs cmd/main.go

EXPOSE 5003
CMD ["./blogs"]