#!/bin/bash

ENDPOINT=${SERVICE_ENDPOINT:-v1}
PORT=${SERVICE_PORT:-80}

CONTAINER_NAME=${SERVICE_NAME:-user-service}
DOCKER_COMPOSE_FILE=${COMPOSE_FILE:-/home/ubuntu/deployments/compose.yaml}

URL="http://localhost:${PORT}/${ENDPOINT}/healthcheck"

echo "Pinging ${URL}..."
# https://stackoverflow.com/questions/38906626/curl-to-return-http-status-code-along-with-the-response
CODE=$(curl -o /dev/null -s -w "%{http_code}" ${URL})

if [[ "${CODE}" -eq 200 ]]; then
  echo "Server is healthy, continuing"
  exit 0
fi

echo "Server returned ${CODE}, attempting to restart"

echo "Stopping existing docker container..."
docker container ls -la --format "{{.Names}}" | \
  grep ${CONTAINER_NAME} | \
  xargs --no-run-if-empty docker stop ${CONTAINER_NAME}

echo "Restarting docker container..."
docker compose -f ${DOCKER_COMPOSE_FILE} restart ${CONTAINER_NAME}
