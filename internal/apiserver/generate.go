//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --package=apiserver --generate gin -o server.gen.go ../../oapi/pmss-schema.yaml
package apiserver
