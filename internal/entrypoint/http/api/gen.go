//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=../../../../openapi/gen-cfg/type.cfg.yml ../../../../openapi/translation-api.yml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=../../../../openapi/gen-cfg/server.cfg.yml ../../../../openapi/translation-api.yml

package api
