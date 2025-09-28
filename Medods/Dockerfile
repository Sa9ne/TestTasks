FROM golang:1.25

WORKDIR /medods-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o medods-service ./cmd

CMD ["./medods-service"]