
# Provide tag as:
# GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
NODE_PORT ?= 3001
SERVER_ORIGIN ?= "http://localhost:3001"
USER_API_BASE_URL ?= "http://user-service:80/v1/users"
GALACTIC_SOVEREIGN_API_BASE_URL ?= "http://galactic-sovereign-service:80/v1/galactic-sovereign"

user-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--tag totocorpsoftwareinc/user-service:${GIT_COMMIT_HASH} \
		-f build/user-service/Dockerfile \
		.

galactic-sovereign-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--tag totocorpsoftwareinc/galactic-sovereign-service:${GIT_COMMIT_HASH} \
		-f build/galactic-sovereign-service/Dockerfile \
		.

user-dashboard-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--build-arg SERVER_ORIGIN=${SERVER_ORIGIN} \
		--build-arg NODE_PORT=${NODE_PORT} \
		--build-arg API_BASE_URL=${USER_API_BASE_URL} \
		--tag totocorpsoftwareinc/user-dashboard:${GIT_COMMIT_HASH} \
		-f build/user-dashboard/Dockerfile \
		.

galactic-sovereign-frontend-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--build-arg SERVER_ORIGIN=${SERVER_ORIGIN} \
		--build-arg NODE_PORT=${NODE_PORT} \
		--build-arg API_BASE_URL=${GALACTIC_SOVEREIGN_API_BASE_URL} \
		--build-arg USER_API_BASE_URL=${USER_API_BASE_URL} \
		--tag totocorpsoftwareinc/galactic-sovereign-frontend:${GIT_COMMIT_HASH} \
		-f build/galactic-sovereign-frontend/Dockerfile \
		.