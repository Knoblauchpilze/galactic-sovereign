
# Provide tag as:
# GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
ENV_DATABASE_HOST ?= "172.17.0.1"
ENV_DATABASE_PORT ?= 5432
ENV_SERVER_PORT ?= 80
# Provide password as:
# ENV_DATABASE_PASSWORD=password

NODE_PORT ?= 3001
SERVER_ORIGIN ?= "http://localhost:3001"
API_BASE_URL ?= "http://galactic-sovereign-service:80/v1"
USER_API_BASE_URL ?= "http://user-service:80/v1/users"

# https://docs.docker.com/network/drivers/bridge/
DOCKER_NETWORK_BRIDGE_NAME ?= "corp-network"

RESTART_RETRIES_COUNT ?= 5

user-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--build-arg ENV_DATABASE_HOST=${ENV_DATABASE_HOST} \
		--build-arg ENV_DATABASE_PORT=${ENV_DATABASE_PORT} \
		--build-arg ENV_SERVER_PORT=${ENV_SERVER_PORT} \
		--build-arg ENV_DATABASE_PASSWORD='${ENV_DATABASE_PASSWORD}' \
		--tag user-service:${GIT_COMMIT_HASH} \
		-f build/user-service/Dockerfile \
		.

user-service-run:
	docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${ENV_DATABASE_PORT} \
		-p ${ENV_SERVER_PORT}:${ENV_SERVER_PORT} \
		-e ENV_DATABASE_HOST=${ENV_DATABASE_HOST} \
		-e ENV_DATABASE_PORT=${ENV_DATABASE_PORT} \
		-e ENV_SERVER_PORT=${ENV_SERVER_PORT} \
		-e ENV_DATABASE_PASSWORD='${ENV_DATABASE_PASSWORD}' \
		--name user-service \
		user-service:${GIT_COMMIT_HASH}

# https://docs.docker.com/config/containers/start-containers-automatically/
user-service-run-detached:
	sudo docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${ENV_DATABASE_PORT} \
		-p ${ENV_SERVER_PORT}:${ENV_SERVER_PORT} \
		-e ENV_DATABASE_HOST=${ENV_DATABASE_HOST} \
		-e ENV_DATABASE_PORT=${ENV_DATABASE_PORT} \
		-e ENV_SERVER_PORT=${ENV_SERVER_PORT} \
		-e ENV_DATABASE_PASSWORD='${ENV_DATABASE_PASSWORD}' \
		--name user-service \
		-d \
		--restart on-failure:${RESTART_RETRIES_COUNT} \
		user-service:${GIT_COMMIT_HASH}

user-service-stop:
	sudo docker stop user-service
	sudo docker rm user-service

user-service-start: user-service-build user-service-run

user-dashboard-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--build-arg SERVER_ORIGIN=${SERVER_ORIGIN} \
		--build-arg NODE_PORT=${NODE_PORT} \
		--build-arg API_BASE_URL=${USER_API_BASE_URL} \
		--tag user-dashboard:${GIT_COMMIT_HASH} \
		-f build/user-dashboard/Dockerfile \
		.

user-dashboard-run:
	docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${NODE_PORT}:${NODE_PORT} \
		user-dashboard:${GIT_COMMIT_HASH}

user-dashboard-run-detached:
	sudo docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${NODE_PORT}:${NODE_PORT} \
		--name user-dashboard \
		-d \
		--restart on-failure:${RESTART_RETRIES_COUNT} \
		user-dashboard:${GIT_COMMIT_HASH}

user-dashboard-stop:
	sudo docker stop user-dashboard
	sudo docker rm user-dashboard

user-dashboard-start: user-dashboard-build user-dashboard-run

stellar-dominion-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--build-arg SERVER_ORIGIN=${SERVER_ORIGIN} \
		--build-arg NODE_PORT=${NODE_PORT} \
		--build-arg API_BASE_URL=${API_BASE_URL} \
		--build-arg USER_API_BASE_URL=${USER_API_BASE_URL} \
		--tag stellar-dominion:${GIT_COMMIT_HASH} \
		-f build/stellar-dominion/Dockerfile \
		.

stellar-dominion-run:
	docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${NODE_PORT}:${NODE_PORT} \
		stellar-dominion:${GIT_COMMIT_HASH}

stellar-dominion-run-detached:
	sudo docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${NODE_PORT}:${NODE_PORT} \
		--name stellar-dominion \
		-d \
		--restart on-failure:${RESTART_RETRIES_COUNT} \
		stellar-dominion:${GIT_COMMIT_HASH}

stellar-dominion-stop:
	sudo docker stop stellar-dominion
	sudo docker rm stellar-dominion

stellar-dominion-start: stellar-dominion-build stellar-dominion-run
