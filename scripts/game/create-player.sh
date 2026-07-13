#!/bin/bash

# This identifier represents the Oberon universe
UNIVERSE="9682f17b-f5f0-4eda-a747-2537d2151837"
API_USER=$(uuidgen)
PLAYER_NAME=""

if [[ $# -ge 1 ]]; then
  PLAYER_NAME=${1}
  echo "Using player name from input: ${PLAYER_NAME}"
else
  # https://stackoverflow.com/questions/1092631/how-to-get-the-current-time-in-seconds-since-the-epoch-in-bash-on-linux
  # https://man7.org/linux/man-pages/man1/date.1.html
  PLAYER_NAME="toto-$(date +%F-%T)"
  echo "No player name provided, using ${PLAYER_NAME}"
fi

PLAYER_OUTPUT_FILE="/tmp/${PLAYER_NAME}.json"
PLANET_OUTPUT_FILE="/tmp/${PLAYER_NAME}_homeworld.json"

BODY="{\"name\":\"${PLAYER_NAME}\",\"api_user\":\"${API_USER}\",\"universe\":\"${UNIVERSE}\"}"

curl -sH 'Content-Type: application/json' \
  http://localhost:60002/v1/galactic-sovereign/players \
  -d ${BODY} \
  -o ${PLAYER_OUTPUT_FILE}

STATUS=$(jq -r '.status' ${PLAYER_OUTPUT_FILE})

if [ "${STATUS}" = "ERROR" ]; then
  echo "Failed to create player:"
  cat ${PLAYER_OUTPUT_FILE}
  echo ""
  rm ${PLAYER_OUTPUT_FILE}
  exit 1
fi

PLAYER_ID=$(jq -r '.details.id' ${PLAYER_OUTPUT_FILE})
HOMEWORLD=$(jq -r '.details.homeworld' ${PLAYER_OUTPUT_FILE})

curl -s \
  http://localhost:60002/v1/galactic-sovereign/planets/${HOMEWORLD} \
  -o ${PLANET_OUTPUT_FILE}

STATUS=$(jq -r '.status' ${PLANET_OUTPUT_FILE})

if [ "${STATUS}" = "ERROR" ]; then
  echo "Failed to get homeworld details player:"
  cat ${PLANET_OUTPUT_FILE}
  echo ""
  rm ${PLANET_OUTPUT_FILE}
  exit 1
fi

echo "Created player ${PLAYER_ID}!"
echo "Homeworld: ${HOMEWORLD}"

PLAYER_SAVE_FILE="sandbox/player.json"
jq -r '.details' ${PLAYER_OUTPUT_FILE} > ${PLAYER_SAVE_FILE}

PLANET_SAVE_FILE="sandbox/planet.json"
jq -r '.details' ${PLANET_OUTPUT_FILE} > ${PLANET_SAVE_FILE}
