build:
	GOOS=linux GOARCH=386 go build -o mussum
	docker build . -t mussum:1.0.0

test:
	go test ./...

run: test build
	docker run --rm -e APIKEY=${APIKEY} mussum:1.0.0
