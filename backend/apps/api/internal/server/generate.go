package server

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -generate types,gin,strict-server -package server -o generated_api.go ../../../specs/001-lms-saas-system/contracts/openapi-lms-v1.yaml
