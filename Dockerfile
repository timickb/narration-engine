FROM golang:1.22.3

WORKDIR /engine
COPY . .
RUN go mod tidy
RUN go build -o engine cmd/main.go

EXPOSE 2140
CMD ["./engine", "-cfg", "cmd/config.yaml"]