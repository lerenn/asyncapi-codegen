//go:generate go run ../../../cmd/asyncapi-codegen -p simple -i ../spec/examples/simple.yml -o ./simple/simple.expected.go
//go:generate go run ../../../cmd/asyncapi-codegen -p correlationID -i ../spec/examples/correlation-id.yml -o ./correlation-id/correlation_id.expected.go

package expected
