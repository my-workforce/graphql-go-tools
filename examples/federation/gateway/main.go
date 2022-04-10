package main

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	log "github.com/jensneuse/abstractlogger"
	"go.uber.org/zap"
	"net"
	"net/http"
	Url "net/url"
	"os"
	"strings"

	"github.com/jensneuse/graphql-go-tools/pkg/graphql"
	"github.com/jensneuse/graphql-go-tools/pkg/playground"

	http2 "github.com/jensneuse/graphql-go-tools/examples/federation/gateway/http"
)

// It's just a simple example of graphql federation gateway server, it's NOT a production ready code.
//
func logger() log.Logger {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}

	return log.NewZapLogger(logger, log.DebugLevel)
}

func startServer() {
	logger := logger()
	logger.Info("logger initialized")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	upgrader := &ws.DefaultHTTPUpgrader
	upgrader.Header = http.Header{}
	upgrader.Header.Add("Sec-Websocket-Protocol", "graphql-ws")

	graphqlEndpoint := "/query"
	playgroundURLPrefix := "/playground"
	playgroundURL := ""

	httpClient := http.DefaultClient

	mux := http.NewServeMux()

	datasourceWatcher := NewDatasourcePoller(httpClient, DatasourcePollerConfig{
		Services: getServicesUrlsFromEnv(),
		//PollingInterval: 30 * time.Second,
	})

	p := playground.New(playground.Config{
		PathPrefix:                      "",
		PlaygroundPath:                  playgroundURLPrefix,
		GraphqlEndpointPath:             graphqlEndpoint,
		GraphQLSubscriptionEndpointPath: graphqlEndpoint,
	})

	handlers, err := p.Handlers()
	if err != nil {
		logger.Fatal("configure handlers", log.Error(err))
		return
	}

	for i := range handlers {
		mux.Handle(handlers[i].Path, handlers[i].Handler)
	}

	var gqlHandlerFactory HandlerFactoryFn = func(schema *graphql.Schema, engine *graphql.ExecutionEngineV2) http.Handler {
		return http2.NewGraphqlHTTPHandler(schema, engine, upgrader, logger)
	}

	gateway := NewGateway(gqlHandlerFactory, httpClient, logger)

	datasourceWatcher.Register(gateway)
	go datasourceWatcher.Run(ctx)

	gateway.Ready()

	mux.Handle("/query", gateway)

	addr := "0.0.0.0:4000"
	logger.Info("Listening",
		log.String("add", addr),
	)
	fmt.Printf("Access Playground on: http://%s%s%s\n", prettyAddr(addr), playgroundURLPrefix, playgroundURL)
	logger.Fatal("failed listening",
		log.Error(http.ListenAndServe(addr, mux)),
	)
}

func prettyAddr(addr string) string {
	return strings.Replace(addr, "0.0.0.0", "localhost", -1)
}

func getServicesUrlsFromEnv() []ServiceConfig {
	var servicesUrls []ServiceConfig
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "URL_") {
			u, err := Url.Parse(pair[1])
			if err != nil {
				log.Error(err)
				continue
			}
			host, _, _ := net.SplitHostPort(u.Host)
			servicesUrls = append(servicesUrls, ServiceConfig{
				Name: host,
				URL:  pair[1],
				WS:   strings.Replace(pair[1], "http", "ws", 1),
			})
		}
	}

	return servicesUrls
}

func main() {
	startServer()
}
