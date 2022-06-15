FROM golang:1.18-bullseye

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go mod download

WORKDIR /app/examples/federation

RUN go mod download

RUN go get github.com/wundergraph/graphql-go-tools/pkg/engine/datasource/graphql_datasource@v1.20.2
RUN go get github.com/wundergraph/graphql-go-tools/pkg/graphqljsonschema@v1.20.2
RUN go get github.com/wundergraph/graphql-go-tools/pkg/engine/datasource/httpclient@v1.20.2

RUN go build -o /tmp/srv-gateway ./gateway

EXPOSE 4000
CMD ["/tmp/srv-gateway"]