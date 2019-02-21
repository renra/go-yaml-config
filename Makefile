SOURCES=./

dep:
	dep init

example:
	go run ${SOURCES}/examples/main.go

.DEFAULT_GOAL := test
test: example
