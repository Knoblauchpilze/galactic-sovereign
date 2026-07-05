#!/bin/bash

# This identifier represents the metal mine
BUILDING="d176e82d-f2ca-4611-996b-c4804096caef"
PLANET_ID=""

if [[ $# -ge 1 ]]; then
  PLANET_ID=${1}
  echo "Using planet id from input: ${PLANET_ID}"
else
  SAVED_PLAYER_FILE="sandbox/player.json"

  if [ ! -f "${SAVED_PLAYER_FILE}" ]; then
    echo "Usage: ./create-building-action.sh player-id"
    echo "Alternatively you can call ./create-player.sh first"
    exit 1
  fi

  PLANET_ID=$(jq -r '.homeworld' ${SAVED_PLAYER_FILE})
  echo "Using player planet (${PLANET_ID}) from file"
fi

OUTPUT_FILE="/tmp/${PLANET_ID}.json"

BODY="{\"building\":\"${BUILDING}\"}"

curl -sH 'Content-Type: application/json' \
  http://localhost:60002/v1/galactic-sovereign/planets/${PLANET_ID}/actions \
  -d ${BODY} \
  -o ${OUTPUT_FILE}

STATUS=$(jq -r '.status' ${OUTPUT_FILE})

if [ "${STATUS}" = "ERROR" ]; then
  echo "Failed to create building action:"
  cat ${OUTPUT_FILE}
  echo ""
  rm ${OUTPUT_FILE}
  exit 1
fi

ACTION_ID=$(jq -r '.details.id' ${OUTPUT_FILE})
COMPLETION_TIME=$(jq -r '.details.completed_at' ${OUTPUT_FILE})

echo "Created building action ${ACTION_ID}!"
echo "Completion time: ${COMPLETION_TIME}"

SAVE_FILE="sandbox/building_action.json"
jq -r '.details' ${OUTPUT_FILE} > ${SAVE_FILE}
