SCHEMA_GEN_DIR=proto/gen

buf.lint:
	buf lint

buf.generate:
	buf generate

run:
	go mod tidy
	go run cmd/main.go --cfg=cmd/config.yaml