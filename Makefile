.PHONY: run
run:
	go run *.go

.PHONY: watch
watch:
	watchman-make -p '*.go' -r 'make run'
