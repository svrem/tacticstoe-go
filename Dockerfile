FROM golang:1.23-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server

CMD ["./server"]