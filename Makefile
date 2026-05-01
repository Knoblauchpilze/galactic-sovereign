
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
SWAG_VERSION ?= v2.0.0-rc5

galactic-sovereign-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--tag totocorpsoftwareinc/galactic-sovereign-service:${GIT_COMMIT_HASH} \
		-f build/galactic-sovereign-service/Dockerfile \
		.

generate-api-spec:
	cd cmd/galactic-sovereign && \
	go run github.com/swaggo/swag/v2/cmd/swag@${SWAG_VERSION} init \
		--v3.1 \
		--generalInfo main.go \
		--dir .,../../internal/controller,../../pkg/communication \
		--output ../../api \
		--outputTypes go,yaml \
		--parseDependencyLevel 1 \
		--parseInternal \
		--generatedTime=false

publish-release:
	./scripts/create-release.sh
