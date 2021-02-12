.PHONY: all clean

all: ugit
	go build

clean: 
	rm ugit && \
	rm -r -f .ugit