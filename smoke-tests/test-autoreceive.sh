#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/util.sh

# --

RECEIVER_TAG=foobar
${NASEFA} auto-receive ${TEMPDIR} ${RECEIVER_TAG} some-tag some-other-tag &
RECEIVE_PID=${!}
log_debug "Started auto-receiver, PID: ${RECEIVE_PID}"

BUNDLE_NAME="send-autoreceive-test-$((RANDOM))"
F1=theplague.png
F2=README.md
F3=LICENSE

${NASEFA} send -bundleName ${BUNDLE_NAME} -to ${RECEIVER_TAG} -to blah ${NASEFA_ROOT}/${F1} ${NASEFA_ROOT}/${F2} ${NASEFA_ROOT}/${F3}
log_debug "Files sent: ${F1} ${F2} ${F3}"

# Give the receiver some time
sleep 2

kill ${RECEIVE_PID}
wait ${RECEIVE_PID} || echo "Autoreceive shut down"

check_files ${NASEFA_ROOT} ${TEMPDIR} ${F1} ${F2} ${F3}
