#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/util.sh

# --

BUNDLE_NAME="web-send-receive-test-whitespace-$((RANDOM))"
F1=README.md
F1_WITH_SPACES="READ ME.md"
F1_URL_ENCODED="READ%20ME.md"

# Create a file with spaces in name
cp ${NASEFA_ROOT}/${F1} "${TEMPDIR}/${F1_WITH_SPACES}"


PORT=8081

${NASEFA} web -allowBundlesListing -bindAddr ":${PORT}" &
WEB_PID=${!}
log_debug "Started web, PID: ${WEB_PID}"
trap "kill -9 ${WEB_PID} || echo" EXIT

${NASEFA} create -bundleName ${BUNDLE_NAME}
log_debug "Bundle created: ${BUNDLE_NAME}"

curl -s -F "file1=@${TEMPDIR}/${F1_WITH_SPACES}" http://localhost:${PORT}/upload/${BUNDLE_NAME} >/dev/null || fail "Failed to upload files"
log_debug "Uploaded 1 file with spaces"

curl -s  "http://localhost:${PORT}/bundle/${BUNDLE_NAME}/${F1_URL_ENCODED}" > ${TEMPDIR}/${F1} || fail "Failed to download file"
log_debug "Downloaded 1 file (saved without spaces)"

kill ${WEB_PID}
wait ${WEB_PID} || echo "Web shut down"

check_files ${NASEFA_ROOT} ${TEMPDIR} ${F1}
