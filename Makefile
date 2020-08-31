clean:
	rm -rf "out"

out:
	mkdir "out"

format:
	go fmt ./...

build: out
	go build -o out/ugo cmd/main.go

test: test-code test-readme test-examples

test-code:
	go test -v ./...

test-readme:
	go run cmd/ugo/main.go run -p ./README.md

test-examples:
	go run cmd/ugo/main.go run -p ./docs/examples/