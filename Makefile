test:
	go test -v ./... -cover -race
	go vet

	staticcheck ./...
