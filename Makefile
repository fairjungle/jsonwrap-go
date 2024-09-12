M = $(shell printf "\033[34;1m▶\033[0m")

all: build

build: ; $(info $(M) Building)
	go build

fmt: ; $(info $(M) Formatting code)
	gofmt -s -w .

full: build tidy

tidy: ; $(info $(M) Tidying mod dep)
	go mod tidy
