generate:
	@echo "Generating code from proto files..."
	@buf generate
	@buf generate --template buf.gen.chatbot.yaml
	@echo "✓ Proto code generation complete"

format:
	gofmt -s -w .

clean:
	@go clean
	@rm -rf ./bin

build: clean
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/app cmd/altalune/*.go

# Utility binaries
build-utils: clean
	@echo "Building utility binaries..."
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/publicid cmd/public_id/*.go
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/secret_encrypter cmd/secret_encrypter/*.go
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/client_secret_hasher cmd/client_secret_hasher/*.go
	@echo "✓ Utility binaries built in ./bin/"

# Individual utility targets
publicid:
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/publicid cmd/public_id/*.go

secret-encrypter:
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/secret_encrypter cmd/secret_encrypter/*.go

client-secret-hasher:
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/client_secret_hasher cmd/client_secret_hasher/*.go
