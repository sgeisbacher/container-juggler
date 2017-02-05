get-deps:
	govendor sync
test:
	go test -v $(shell go list ./... | grep -v /vendor/)
