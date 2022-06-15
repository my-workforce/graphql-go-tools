FROM golang:1.18-bullseye

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download

WORKDIR /app/examples/federation

RUN go mod download

RUN go build -o /tmp/srv-gateway ./gateway

EXPOSE 4000
CMD ["/tmp/srv-gateway"]