//go:generate go run ../../../cmd/asyncapi-codegen -p anyof         -i ../spec/examples/anyof.yml          -o ./anyof/anyof.expected.go
//go:generate go run ../../../cmd/asyncapi-codegen -p correlationID -i ../spec/examples/correlation-id.yml -o ./correlation-id/correlation_id.expected.go
//go:generate go run ../../../cmd/asyncapi-codegen -p oneof         -i ../spec/examples/oneof.yml          -o ./oneof/oneof.expected.go
//go:generate go run ../../../cmd/asyncapi-codegen -p rpcServer     -i ../spec/examples/rpc-server.yml     -o ./rpc-server/rpc-server.expected.go
//go:generate go run ../../../cmd/asyncapi-codegen -p simple        -i ../spec/examples/simple.yml         -o ./simple/simple.expected.go

package expected
