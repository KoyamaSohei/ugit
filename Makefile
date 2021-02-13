.PHONY: all clean

all: ugit

ugit: cmd/*.go data/*.go
	go build

clean: 
	rm ugit && \
	rm -r -f .ugit