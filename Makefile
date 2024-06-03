
# Provide tag as:
# GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
ENV_DATABASE_HOST ?= "172.17.0.1"
ENV_DATABASE_PORT ?= 5432
ENV_SERVER_PORT ?= 80
# Provide password as:
# ENV_DATABASE_PASSWORD=password

NODE_PORT ?= 3001
SERVER_ORIGIN ?= "http://localhost:3001"
API_BASE_URL ?= "http://user-service:80/v1/users"

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
		-f build/users/Dockerfile \
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

webserver-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--build-arg SERVER_ORIGIN=${SERVER_ORIGIN} \
		--build-arg NODE_PORT=${NODE_PORT} \
		--build-arg API_BASE_URL=${API_BASE_URL} \
		--tag webserver:${GIT_COMMIT_HASH} \
		-f build/webserver/Dockerfile \
		.

webserver-run:
	docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${NODE_PORT}:${NODE_PORT} \
		webserver:${GIT_COMMIT_HASH}

webserver-run-detached:
	sudo docker run \
		--network ${DOCKER_NETWORK_BRIDGE_NAME} \
		-p ${NODE_PORT}:${NODE_PORT} \
		--name webserver \
		-d \
		--restart on-failure:${RESTART_RETRIES_COUNT} \
		webserver:${GIT_COMMIT_HASH}

webserver-stop:
	sudo docker stop webserver
	sudo docker rm webserver

webserver-start: webserver-build webserver-run
