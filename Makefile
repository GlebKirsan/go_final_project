.DEFAULT_GOAL = build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build cmd/task-tracker/main.go
.PHONY:build

run: build
	./main
.PHONY:run