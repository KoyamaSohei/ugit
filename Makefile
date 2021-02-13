.PHONY: all clean test

all: ugit

ugit: *.go data/*.go base/*.go
	go build

clean: 
	rm ugit && \
	rm -r -f .ugit

test:
	go test