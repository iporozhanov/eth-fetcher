.PHONY: run build-docker run-docker test
include .env

run:
	go mod tidy && go run .

build-docker:
	docker build --build-arg API_PORT=${API_PORT} -t ethfetcher .

run-docker:
	docker run --env-file .env -p ${API_PORT}:${API_PORT} ethfetcher
test:
	go test -v ./...