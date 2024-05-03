#!/bin/bash

ENDPOINT=${SERVICE_ENDPOINT:-v1}
PORT=${SERVICE_PORT:-80}

CONTAINER_NAME=${SERVICE_NAME:-user-service}

URL="http://localhost:${PORT}/${ENDPOINT}/healthcheck"

echo "Pinging ${URL}..."
CODE=$(curl -o /dev/null -s -w "%{http_code}" ${URL})

if [[ "${CODE}" -eq 200 ]]; then
  echo "Server is healthy, continuing"
  exit 0
fi

echo "Server returned ${CODE}, attempting to restart"

echo "Stopping existing docker container..."
sudo docker container ls -la --format "{{.Names}}" | \
  grep ${CONTAINER_NAME} | \
  xargs --no-run-if-empty sudo docker stop ${CONTAINER_NAME}

echo "Restarting docker container..."
sudo docker start ${CONTAINER_NAME}
