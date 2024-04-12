#!/bin/bash

if [ "$#" -ne 2 ]; then
  echo "usage: ec2-establish-db-tunnel.sh path/to/identify/file ec2-ip-address"
  exit 1
fi

LOCAL_PORT=${LOCAL_DATABASE_PORT:-5000}
REMOTE_PORT=${REMOTE_DATABASE_PORT:-5432}

EC2_USER=${EC2_USER_NAME:-ubuntu}
DB_NAME=${DATABASE_NAME:-db_user_service}
DB_USER=${DATABASE_USER:-user_service_admin}

PATH_TO_SOCKET="/tmp/ec2-db-tunnel"
SOCKET_SERVER_NAME="db-tunnel"

EC2_IP_ADDRESS=$1
SSH_IDENTITY_FILE=$2

if [ ! -f "${SSH_IDENTITY_FILE}" ]; then
  echo "Can't access identify file ${SSH_IDENTITY_FILE}, aborting"
  exit 1
fi

if [ -S "${PATH_TO_SOCKET}" ]; then
  echo "SSH tunnel already in use, please close it first"
  exit 0
fi

echo "Creating tunnel on port ${LOCAL_PORT}..."

# https://linuxize.com/post/how-to-setup-ssh-tunneling/
ssh -i ${SSH_IDENTITY_FILE} \
  -NfM \
  -S ${PATH_TO_SOCKET} \
  -L ${LOCAL_PORT}:localhost:${REMOTE_PORT} \
  ${EC2_USER}@${EC2_IP_ADDRESS} \
  ${SOCKET_SERVER_NAME}
