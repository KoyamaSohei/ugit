.PHONY: all clean test

all: ugit

ugit: cmd/*.go data/*.go
	go build

clean: 
	rm ugit && \
	rm -r -f .ugit

test:
	go test