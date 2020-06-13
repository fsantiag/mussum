build: test lint
	GOOS=linux GOARCH=386 go build -o mussum
	docker build . -t fsantiag/mussum:1.0.0

lint:
	golint ./...

test:
	go test ./...

run: test build
	docker run --rm -e APIKEY=${APIKEY} mussum:1.0.0
