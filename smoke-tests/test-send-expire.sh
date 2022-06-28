#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/util.sh

# --

BUNDLE_NAME="send-expire-test-$((RANDOM))"
F1=theplague.png

${NASEFA} send -expire 1s -bundleName ${BUNDLE_NAME} ${NASEFA_ROOT}/${F1}
log_debug "Files sent: ${F1}"

sleep 2

${NASEFA} receive ${TEMPDIR} ${BUNDLE_NAME}

check_no_files ${TEMPDIR}
