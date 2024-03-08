#!/bin/bash
bold=$(tput bold)
normal=$(tput sgr0)

# Usage function
display_usage () {
  echo "${bold}Run the tests in a specific directory.${normal}"
  echo "Usage: ./test.sh directory"
  echo "${bold}Available options:${normal}"
  echo "--dir      : directory to test"
  echo "--usage    : displays this usage"
  echo "-h --help  : displays this usage"
}

# General options.
DISPLAY_USAGE="false"

TEST_DIRECTORY=""

parse_general_config () {
  ARGUMENTS=()
  while [[ $# -gt 0 ]]
  do
    key="$1"

    case $key in
      --dir)
        TEST_DIRECTORY=${2}
        shift # past argument
        shift # past value
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
if [ "${DISPLAY_USAGE}" = "true" ] || [ "${TEST_DIRECTORY}" = "" ]; then
  display_usage
  exit 1
fi

# Check if go tests exist here.
GO_TEST_FILES_COUNT=$(find ${TEST_DIRECTORY} -type f -name "*_test.go" | wc -l)

if [ ${GO_TEST_FILES_COUNT} -eq 0 ]; then
  echo "No tests found in ${TEST_DIRECTORY}"
  exit 0
fi

# Collect sub-folders recursively (including .).
cd ${TEST_DIRECTORY}
echo "Testing ${TEST_DIRECTORY}"
go test ./...

# Register the exit code and propagate it if it is
# non-zero (to indicate failure of the sub-program).
EXIT_CODE=$?
if [ ${EXIT_CODE} -ne 0 ]; then
  echo "Testing directory ${TEST_DIRECTORT} returned exit code ${EXIT_CODE}"
  exit 1
fi
