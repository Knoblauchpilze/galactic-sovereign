#!/bin/bash

ENDPOINT=${SERVICE_ENDPOINT:-v1/users}

CONTAINER_NAME=${SERVICE_NAME:-user-service}
DOCKER_COMPOSE_FILE=${COMPOSE_FILE:-/home/ubuntu/deployments/compose.yaml}

# https://stackoverflow.com/questions/39070547/how-to-expose-a-docker-network-to-the-host-machine
echo "Fetching IP of ${CONTAINER_NAME}"
CONTAINER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${CONTAINER_NAME})

URL="http://${CONTAINER_IP}/${ENDPOINT}/healthcheck"

echo "Pinging ${URL}..."
# https://stackoverflow.com/questions/38906626/curl-to-return-http-status-code-along-with-the-response
CODE=$(curl -o /dev/null -s -w "%{http_code}" ${URL})

if [[ "${CODE}" -eq 200 ]]; then
  echo "Server is healthy, continuing"
  exit 0
fi

echo "Server returned ${CODE}, attempting to restart"

echo "Stopping existing docker container..."
docker container ls --format "{{.Names}}" | \
  grep ${CONTAINER_NAME} | \
  xargs --no-run-if-empty docker stop ${CONTAINER_NAME}

echo "Restarting docker container..."
docker compose -f ${DOCKER_COMPOSE_FILE} restart ${CONTAINER_NAME}
