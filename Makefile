.PHONY: all build install clean

build:
	go build -ldflags="-w -s"

install:
	mv zoomeye-cli /usr/bin/zoomeye-cli

all: build install

clean:
	rm -f /usr/bin/zoomeye-cli
