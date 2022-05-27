#!/usr/bin/env bash

# Script for local smoke testing

set -e

TEST_DIR="./nasefa-test"
NASEFA="./nasefa"

FILES_SIZE=10000
FILE1="${TEST_DIR}/file1.bin"
FILE2="${TEST_DIR}/file2.bin"

rm -rf "${TEST_DIR}"
mkdir -p "${TEST_DIR}"

go build

head -c ${FILES_SIZE} /dev/random > "${FILE1}"
head -c ${FILES_SIZE} /dev/random > "${FILE2}"

${NASEFA} send --fileId f1 ${FILE1}
${NASEFA} send --fileId f2 ${FILE2}

mv "${FILE1}" "${FILE1}.original"
mv "${FILE2}" "${FILE2}.original"

${NASEFA} list

${NASEFA} receive "${TEST_DIR}" f1
${NASEFA} receive "${TEST_DIR}" f2

diff -q "${FILE1}.original" "${FILE1}"
diff -q "${FILE2}.original" "${FILE2}"
