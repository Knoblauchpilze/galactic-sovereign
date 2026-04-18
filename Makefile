
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
SWAG_VERSION ?= v1.16.6

galactic-sovereign-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--tag totocorpsoftwareinc/galactic-sovereign-service:${GIT_COMMIT_HASH} \
		-f build/galactic-sovereign-service/Dockerfile \
		.

generate-api-spec:
	cd cmd/galactic-sovereign && \
	go run github.com/swaggo/swag/cmd/swag@${SWAG_VERSION} init \
		--generalInfo main.go \
		--dir .,../../internal/controller \
		--output ../../api \
		--outputTypes yaml \
		--parseInternal \
		--generatedTime=false
