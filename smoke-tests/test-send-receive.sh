#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/util.sh

# --

BUNDLE_NAME="send-receive-test-$((RANDOM))"
F1=theplague.png
F2=README.md
F3=LICENSE

${NASEFA_BIN} send -bundleName ${BUNDLE_NAME} ${NASEFA_ROOT}/${F1} ${NASEFA_ROOT}/${F2} ${NASEFA_ROOT}/${F3}
log_debug "Files sent: ${F1} ${F2} ${F3}"

${NASEFA_BIN} receive ${TEMPDIR} ${BUNDLE_NAME}

check_files ${NASEFA_ROOT} ${TEMPDIR} ${F1} ${F2} ${F3}
