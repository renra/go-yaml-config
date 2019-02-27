SOURCES=./

dep:
	dep init

example:
	go run ${SOURCES}/examples/main.go

.PHONY: test
.DEFAULT_GOAL := test
test:
	go test ./test... -count 1 -v
