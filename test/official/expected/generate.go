//go:generate go run ../../../cmd/asyncapi-codegen -p test -i ../spec/examples/simple.yml -o ./simple/simple.expected.go
//go:generate go run ../../../cmd/asyncapi-codegen -p test -i ../spec/examples/correlation-id.yml -o ./correlation-id/correlation_id.expected.go

package expected
