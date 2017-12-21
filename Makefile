.PHONY: all build

all: build

build:
	go build -o build/spotify ./cmd/spotify/...
