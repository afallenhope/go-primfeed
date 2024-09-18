build:
	@go build -o bin/primfeed cmd/main.go

run: build
	@./bin/primfeed
