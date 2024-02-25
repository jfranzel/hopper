build:
	go build -o bin/hopper cmd/main.go
run:
	./bin/hopper
test:
	go test -v ./...