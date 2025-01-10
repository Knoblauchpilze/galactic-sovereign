
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)

galactic-sovereign-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--tag totocorpsoftwareinc/galactic-sovereign-service:${GIT_COMMIT_HASH} \
		-f build/galactic-sovereign-service/Dockerfile \
		.