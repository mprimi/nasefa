#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/util.sh

# --

BUNDLE_NAME="web-send-receive-test-$((RANDOM))"
F1=theplague.png
F2=README.md

PORT=8081

${NASEFA_BIN} web -bindAddr ":${PORT}" &
WEB_PID=${!}
log_debug "Started web, PID: ${WEB_PID}"

${NASEFA_BIN} create -bundleName ${BUNDLE_NAME}
log_debug "Bundle created: ${BUNDLE_NAME}"

curl -s http://localhost:${PORT}/bundle/${BUNDLE_NAME} >/dev/null || fail "Failed to list bundle ${BUNDLE_NAME}"
log_debug "Bundle listed"

curl -s -F "file1=@${NASEFA_ROOT}/${F1}" -F "file2=@${NASEFA_ROOT}/${F2}" http://localhost:${PORT}/upload/${BUNDLE_NAME} >/dev/null || fail "Failed to upload files"
log_debug "Uploaded 2 files"

curl -s  http://localhost:${PORT}/bundle/${BUNDLE_NAME}/${F1} > ${TEMPDIR}/${F1} || fail "Failed to download file"
curl -s  http://localhost:${PORT}/bundle/${BUNDLE_NAME}/${F2} > ${TEMPDIR}/${F2} || fail "Failed to download file"
log_debug "Downloaded 2 files"

kill ${WEB_PID}
wait ${WEB_PID} || echo "Web shut down"

check_files ${NASEFA_ROOT} ${TEMPDIR} ${F1} ${F2}
