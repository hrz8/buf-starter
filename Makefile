format:
	gofmt -s -w .

clean:
	@go clean
	@rm -rf ./bin

build: clean
	@env GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/app cmd/cli/*.go
