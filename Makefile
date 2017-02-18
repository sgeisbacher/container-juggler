get-deps:
	govendor sync
test:
	go test -cover -v $(shell go list ./... | grep -v /vendor/)
build:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/container-juggler-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/container-juggler-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o ./bin/container-juggler-win64.exe

