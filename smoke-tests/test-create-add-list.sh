#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/util.sh

# --

BUNDLE_NAME="create-add-list-test-$((RANDOM))"
F1=theplague.png
F2=README.md
F3=LICENSE

${NASEFA_BIN} create -bundleName ${BUNDLE_NAME}
log_debug "Bundle created ${BUNDLE_NAME}"

${NASEFA_BIN} add -bundleName ${BUNDLE_NAME} ${NASEFA_ROOT}/${F1} ${NASEFA_ROOT}/${F2}
log_debug "Files added: ${F1} ${F2}"

${NASEFA_BIN} list
# TODO check listing
