//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --package=apiclient --generate client,types -o client.gen.go ../../oapi/pmss-schema.yaml
package apiclient
