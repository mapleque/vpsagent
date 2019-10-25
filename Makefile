.PHONY: all build linux mac

build: bin
	go build -o bin/vpsagent

linux: bin/vpsagent-linux

mac: bin/vpsagent-mac

all: bin/vpsagent-linux bin/vpsagent-mac

bin/vpsagent-linux: bin
	GOOS=linux GOARCH=amd64 go build -o bin/vpsagent-linux main.go
	@echo '>bin/vpsagent-linux'

bin/vpsagent-darwin: bin
	GOOS=darwin GOARCH=amd64 go build -o bin/vpsagent-mac main.go
	@echo '>bin/vpsagent-mac'

bin:
	mkdir bin

clean:
	rm -rf bin
