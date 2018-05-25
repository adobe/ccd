.PHONY: deps coverage

deps:
	go get -u github.com/AlekSi/gocoverutil

coverage:
	gocoverutil -coverprofile=coverage.out test -v -covermode=count github.com/adobe/ccd/...
	go tool cover -html=coverage.out

test:
	go test -v -tags all ./...
integration-test:
	go test -v -tags integration-test ./...
