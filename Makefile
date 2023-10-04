.PHONY: build

build:
	go build -v

image:
	docker build -t tg-link-bot:v0.1 .

container:
	docker run --name tg-link-bot -p 4000:4000 -d tg-link-bot:v0.1

.DEFAULT_GOAL := build
