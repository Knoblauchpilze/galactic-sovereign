#!/bin/bash
bold=$(tput bold)
normal=$(tput sgr0)

# Usage function
display_usage () {
  echo "${bold}Run the specified application with provided configuration.${normal}"
  echo "Usage: ./run.sh --app-name name_of_app"
  echo "${bold}Available options:${normal}"
  echo "--usage           : displays this usage"
  echo "-h --help         : displays this usage"
  echo "--app-name        : the name of the app to execute"
  echo "--app-path        : the path to reach the executable of the app to run"
  echo "--ld-library-path : additional library paths to set for the execution"
  echo "--debug           : run the app in debug mode"
}

# General options.
DISPLAY_USAGE="false"

APP_NAME=""
APP_PATH=""

ADDITIONAL_LD_LIBRARY_PATH="/usr/local/lib/:/usr/lib/x86_64-linux-gnu"

DEBUG_MODE="false"

parse_general_config () {
  ARGUMENTS=()
  while [[ $# -gt 0 ]]
  do
    key="$1"

    case $key in
      --app-name)
        APP_NAME=${2}
        shift # past argument
        shift # past value
      ;;
      --app-path)
        APP_PATH=${2}
        shift # past argument
        shift # past value
      ;;
      --ld-library-path)
        ADDITIONAL_LD_LIBRARY_PATH=${2}
        shift # past argument
        shift # past value
      ;;
      --debug)
        DEBUG_MODE="true"
        shift # past argument
      ;;
      -h|--help|--usage)
        DISPLAY_USAGE="true"
        shift # past argument
      ;;
      *) # unknown option
        ARGUMENTS+=("$1") # save it in an array for later
        shift # past argument
      ;;
    esac
  done
  set -- "${ARGUMENTS[@]}"
}

# Parse input arguments.
ARGS=()
while [[ $# -gt 0 ]]; do
  ARGS+="$1 "
  shift # past argument
done

parse_general_config $ARGS

# Exit early if usage is required or in case
# a missing value is detected.
if [ "${DISPLAY_USAGE}" = "true" ] || [ "${APP_NAME}" = "" ]; then
  display_usage
  exit 1
fi

# Generate processed arguments.
FULL_APP_FILE="./${APP_NAME}"
if [ "${APP_PATH}" != "" ]; then
  FULL_APP_FILE="${APP_PATH}/${APP_NAME}"
fi

# Export additional library paths.
export LD_LIBRARY_PATH=$PWD/lib:${ADDITIONAL_LD_LIBRARY_PATH}

# Execute the app.
if [ "${DEBUG_MODE}" == "true" ]; then
  gdb --args ${FULL_APP_FILE}
else
  ${FULL_APP_FILE}
fi

# Register the exit code and propagate it if it is
# non-zero (to indicate failure of the sub-program).
EXIT_CODE=$?
if [ ${EXIT_CODE} -ne 0 ]; then
  echo "Application ${APP_NAME} exited with status ${EXIT_CODE}"
  exit 1
fi
