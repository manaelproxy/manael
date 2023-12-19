.PHONY: all
all: bin/manael

bin/%: cmd/%/main.go *.go internal/**/*.go
	@mkdir -p bin
	go build -o $@ $<

.PHONY: clean
clean:
	rm -r bin
