all: native

bin:
	mkdir -p $@

native: bin
	go build -o bin/native cmd/native/main.go
