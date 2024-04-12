#!/bin/bash

PATH_TO_SOCKET="/tmp/ec2-db-tunnel"
SOCKET_SERVER_NAME="db-tunnel"

if [ ! -S "${PATH_TO_SOCKET}" ]; then
  echo "It seems no ssh tunnel is up and running, exiting"
  exit 0
fi

# https://unix.stackexchange.com/questions/83806/how-to-kill-ssh-session-that-was-started-with-the-f-option-run-in-background
ssh -S ${PATH_TO_SOCKET} -O exit ${SOCKET_SERVER_NAME}
