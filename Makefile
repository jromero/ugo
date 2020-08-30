clean:
	rm -rf "out"

out:
	mkdir "out"

test:
	go test -v ./...

build: out
	go build -o out/ugo cmd/main.go

test-examples:
	go run cmd/main.go run -p ./docs/examples/