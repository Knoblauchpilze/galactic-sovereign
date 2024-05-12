#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "usage: website-run.sh path/to/web/server/artifacts"
  exit 1
fi

SERVER_PATH=${1}
# https://stackoverflow.com/questions/13333221/how-to-change-value-of-process-env-port-in-node-js
SERVER_PORT=${NODE_PORT:-3000}
SERVER_URL="http://localhost:${SERVER_PORT}"

echo "Starting server ${SERVER_PATH} at ${SERVER_URL}..."
ORIGIN=${SERVER_URL} PORT=${SERVER_PORT} node ${SERVER_PATH}

if [ $? != 0 ]; then
  echo "Webserver crashed, exiting"
  exit 1
fi

echo "Webserver shutdown gracefully"
