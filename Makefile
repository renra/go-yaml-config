SOURCES=./

dep:
	dep init

non_panicking_example:
	go run ${SOURCES}/examples/non_panicking/main.go

panicking_example:
	go run ${SOURCES}/examples/panicking/main.go

examples: non_panicking_example panicking_example

.PHONY: test
.DEFAULT_GOAL := test
test:
	go test ./test... -count 1 -v
