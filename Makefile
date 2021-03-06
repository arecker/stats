.PHONY: all
all: bin/run

bin/run: $(wildcard src/*.go)
	go build -o $@ ./...

.PHONY: run
run: bin/run
	./bin/run

.PHONY: test
test:
	go test ./...

.PHONY: update
update:
	git submodule update --init --recursive --remote

.PHONY: clean
clean:
	rm bin/run
	rm words/*.json
