#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/util.sh

# --

BUNDLE_NAME="send-receive-test-$((RANDOM))"
F1=theplague.png
F2=README.md
F3=LICENSE

${NASEFA} send -bundleName ${BUNDLE_NAME} ${NASEFA_ROOT}/${F1} ${NASEFA_ROOT}/${F2} ${NASEFA_ROOT}/${F3}
log_debug "Files sent: ${F1} ${F2} ${F3}"

${NASEFA} delete -bundleName ${BUNDLE_NAME}

${NASEFA} receive ${TEMPDIR} ${BUNDLE_NAME} || is_error=true

if [[ "${is_error}" != "true" ]]; then
  fail "Expected error downloading deleted bundle"
fi

echo "Downloading deleted bundle failed (as expected)"
