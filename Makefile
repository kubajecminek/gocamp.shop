all: build test
 
build:
	mkdir -p out/bin
	go build -o out/bin/gocamp *.go

run:
	go run *.go
 
test:
	go test ./...

bench:
	go test -bench=.
 
clean:
	go clean
	rm out/bin/gocamp
