test:
	go test -v ./... -cover -race
	go vet

	go get -u honnef.co/go/tools/cmd/staticcheck

	staticcheck ./...
