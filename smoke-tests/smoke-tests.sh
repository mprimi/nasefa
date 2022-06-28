#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd ) # https://stackoverflow.com/a/246128

NASEFA_ROOT="${SCRIPT_DIR}/.."

tests=$(ls -1 ${SCRIPT_DIR}/test-*.sh)
failed_tests=""

for t in ${tests};
do
  echo
  echo '**'
  echo "** TEST: ${t}"
  echo '**'
  echo
  bash ${t} || failed_tests="${failed_tests} `basename $t`"
done

echo
echo "---"
echo

if [[ "${failed_tests}" != "" ]]; then
  echo "‼️  Failed tests: ${failed_tests}"
  exit 1
fi


echo "✅ All tests passed"
