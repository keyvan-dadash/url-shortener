FROM golang AS builder

WORKDIR /go/src/url

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o url .

FROM ubuntu:latest  

WORKDIR /root/

COPY --from=builder /go/src/url .

EXPOSE 8080

CMD ["./url"]