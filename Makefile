
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
SWAG_VERSION ?= v2.0.0-rc5
GOLANGCI_LINT_VERSION ?= v2.12.2

galactic-sovereign-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--tag totocorpsoftwareinc/galactic-sovereign-service:${GIT_COMMIT_HASH} \
		-f build/galactic-sovereign-service/Dockerfile \
		.

generate-mocks:
	go generate ./...

generate-api-spec:
	cd cmd/galactic-sovereign && \
	go run github.com/swaggo/swag/v2/cmd/swag@${SWAG_VERSION} init \
		--v3.1 \
		--generalInfo main.go \
		--dir .,../../pkg/domain/adapters/driving,../../pkg/domain/adapters/driving/dtos \
		--output ../../api \
		--outputTypes go,yaml \
		--parseDependencyLevel 1 \
		--parseInternal \
		--generatedTime=false

publish-release:
	./scripts/create-release.sh

tests:
	go test ./...

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} run ./...

fix-lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} run --fix ./...