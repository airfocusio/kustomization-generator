.PHONY: *

run: run-download run-helm run-kustomize

run-download:
	go run . --dir=example/download

run-helm:
	go run . --dir=example/helm

run-kustomize:
	go run . --dir=example/kustomize

test:
	go test -v ./...

test-watch:
	watch -n1 go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build:
	goreleaser release --rm-dist --skip-publish --snapshot

release:
	goreleaser release --rm-dist
