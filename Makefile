all: native

bin:
	mkdir -p $@

native: bin/native-agent

.PHONY: bin/native-agent
bin/native-agent: bin
	CGO_ENABLED=0 go build -o $@ cmd/native/main.go
