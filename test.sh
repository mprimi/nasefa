#!/usr/bin/env bash

# Script for local smoke testing

set -e

TEST_DIR="./nasefa-test"
NASEFA="./nasefa"

FILES_SIZE=10000
FILE1="${TEST_DIR}/document.pdf"
FILE2="${TEST_DIR}/image.jpg"
FILE3="${TEST_DIR}/archive.zip"

rm -rf "${TEST_DIR}"
mkdir -p "${TEST_DIR}"

go build

head -c ${FILES_SIZE} /dev/random > "${FILE1}"
head -c ${FILES_SIZE} /dev/random > "${FILE2}"
head -c ${FILES_SIZE} /dev/random > "${FILE3}"

${NASEFA} send -bundleName documents ${FILE1} ${FILE2}
${NASEFA} send -bundleName archive -to foo -to bar ${FILE3}
${NASEFA} send -bundleName archive-autodelete -expire "1s" ${FILE3}

mv "${FILE1}" "${FILE1}.original"
mv "${FILE2}" "${FILE2}.original"
mv "${FILE3}" "${FILE3}.original"

${NASEFA} list

${NASEFA} receive "${TEST_DIR}" documents
${NASEFA} receive "${TEST_DIR}" archive

diff -q "${FILE1}.original" "${FILE1}"
diff -q "${FILE2}.original" "${FILE2}"
diff -q "${FILE3}.original" "${FILE3}"


sleep 2
# This is expected to fail!
${NASEFA} receive "${TEST_DIR}" archive-autodelete
