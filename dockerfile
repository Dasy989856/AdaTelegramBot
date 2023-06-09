FROM golang:alpine
WORKDIR /ada
COPY . .
RUN go build -o adaApp ./cmd/main.go
CMD ["./adaApp"]