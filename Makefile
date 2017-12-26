.PHONY: all build

all: build

build:
	go build -o build/spotctl ./cmd/spotctl/...
