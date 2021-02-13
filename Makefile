.PHONY: all clean test

all: ugit

ugit: data/*.go
	go build

clean: 
	rm ugit && \
	rm -r -f .ugit

test:
	go test