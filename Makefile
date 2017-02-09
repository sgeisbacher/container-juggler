get-deps:
	govendor sync
test:
	go test -cover -v $(shell go list ./... | grep -v /vendor/)
