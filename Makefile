.PHONY: default

default: di build;

di:
	wire ./cmd/user

build:
	go build -o out/user ./cmd/user