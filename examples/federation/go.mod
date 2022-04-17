module github.com/jensneuse/graphql-go-tools/examples/federation

go 1.15

require (
	github.com/gobwas/ws v1.1.0
	github.com/jensneuse/abstractlogger v0.0.4
	github.com/jensneuse/graphql-go-tools v1.20.2
	github.com/nats-io/nats-server/v2 v2.3.2 // indirect
	go.uber.org/zap v1.18.1
	"github.com/rs/cors" v1.7.0
)

replace github.com/jensneuse/graphql-go-tools => ../../
