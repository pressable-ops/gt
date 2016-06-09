all:
	go build
	env GOOS=darwin GOARCH=amd64 go build -o gtsign-darwin-amd64
	env GOOS=linux GOARCH=amd64 go build -o gtsign-linux-amd64
	env GOOS=linux GOARCH=386 go build -o gtsign-linux-386
