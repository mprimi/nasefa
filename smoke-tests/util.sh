#!/usr/bin/env bash
NASEFA_ROOT="${SCRIPT_DIR}/.."
NASEFA_BIN="${NASEFA_ROOT}/nasefa"

function fail() {
  echo "❌ $*"
  exit 1
}

function check_files() {
  local source_dir=${1}
  local target_dir=${2}
  shift 2;
  for fname in ${@};
  do
    local source_file="${source_dir}/${fname}"
    local target_file="${target_dir}/${fname}"
    log_debug "Checking file ${fname}"
    diff -q "${source_file}" "${target_file}" || fail "File ${F1} does not match!"
  done
}

function check_no_files() {
  local dir_list=`ls -1 ${1}`

  if [[ "${dir_list}" != "" ]]; then
    fail "Expecting no files, but found: ${dir_list}"
  fi
}

function log_debug() {
  echo "🐛 $*"
}

test -d "${NASEFA_ROOT}" || fail "Nasefa root directory not found at: ${NASEFA_ROOT}"
log_debug "Nasefa project root: ${NASEFA_ROOT}"

test -x "${NASEFA_BIN}" || fail "Nasefa binary not found at: ${NASEFA_BIN}"
log_debug "Nasefa binary: ${NASEFA_BIN}"

TEMPDIR=$(mktemp -q -d) || fail "Could not create temporary folder"
log_debug "Temporary folder: ${TEMPDIR}"
