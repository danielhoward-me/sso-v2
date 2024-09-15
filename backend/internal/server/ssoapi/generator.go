package ssoapi

//go:generate redocly bundle -o ./schema.gen.yaml ./../schemas/sso.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config ./../oapi-codegen.yaml -package ssoapi ./schema.gen.yaml
