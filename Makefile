
# Provide tag as:
# GIT_COMMIT_HASH=$(git rev-parse --short HEAD)
ENV_DATABASE_HOST="172.17.0.1"
ENV_DATABASE_PORT=5432
ENV_SERVER_PORT=80
# Provide password as:
# ENV_DATABASE_PASSWORD=password

RESTART_RETRIES_COUNT=5

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
		-p ${ENV_DATABASE_PORT} \
		-p ${ENV_SERVER_PORT}:${ENV_SERVER_PORT} \
		-e ENV_DATABASE_HOST=${ENV_DATABASE_HOST} \
		-e ENV_DATABASE_PORT=${ENV_DATABASE_PORT} \
		-e ENV_SERVER_PORT=${ENV_SERVER_PORT} \
		-e ENV_DATABASE_PASSWORD='${ENV_DATABASE_PASSWORD}' \
		user-service:${GIT_COMMIT_HASH}

# https://docs.docker.com/config/containers/start-containers-automatically/
user-service-run-detached:
	sudo docker run \
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
