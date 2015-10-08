# Docker Makefile for rabbit-herder
#
NAME=rabbit-herder
PWD=$$(pwd)

all: builder build clean

build: builder
	docker build -t micahhausler/$(NAME) .

builder:
	docker build -t micahhausler/$(NAME)-builder -f Dockerfile.build .
	docker run --rm -v $(PWD):/usr/src/app micahhausler/$(NAME)-builder

clean:
	rm rabbit-herder
