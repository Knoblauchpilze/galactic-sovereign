#!/bin/bash

PLAYER_ID=""

if [[ $# -ge 1 ]]; then
  PLAYER_ID=${1}
  echo "Using player id from input: ${PLAYER_ID}"
else
  SAVED_PLAYER_FILE="sandbox/player.json"

  if [ ! -f "${SAVED_PLAYER_FILE}" ]; then
    echo "Usage: ./create-building-action.sh player-id"
    echo "Alternatively you can call ./create-player.sh first"
    exit 1
  fi

  PLAYER_ID=$(jq -r '.id' ${SAVED_PLAYER_FILE})
  echo "Using player from file ${PLAYER_ID}"
fi

OUTPUT_FILE="/tmp/${PLAYER_NAME}.json"p

curl -sH 'Content-Type: application/json' \
  -X POST \
  http://localhost:60002/v1/galactic-sovereign/players/${PLAYER_ID}/planets \
  -o ${OUTPUT_FILE}

STATUS=$(jq -r '.status' ${OUTPUT_FILE})

if [ "${STATUS}" = "ERROR" ]; then
  echo "Failed to create planet:"
  cat ${OUTPUT_FILE}
  echo ""
  rm ${OUTPUT_FILE}
  exit 1
fi

PLANET_ID=$(jq -r '.details.id' ${OUTPUT_FILE})

echo "Created planet ${PLANET_ID}!"

SAVE_FILE="sandbox/planet.json"
jq -r '.details' ${OUTPUT_FILE} > ${SAVE_FILE}
