package ssointernalapi

//go:generate redocly bundle -o ./schema.gen.yaml ./../schemas/sso-internal.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config ./../oapi-codegen.yaml -package ssointernalapi ./schema.gen.yaml
