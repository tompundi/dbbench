all:
	go build -o kvbench cmd/kvbench/main.go

test:
	go test -v .