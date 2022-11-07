.PHONY: run
run:
	go run *.go

.PHONY: watch
watch:
	watchman-make -p '**/*.go' -r 'make run'

.PHONY: build
build:
	go build

.PHONY: start-for-debug
start-for-debug: build
	./retool
