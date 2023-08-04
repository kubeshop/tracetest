#!/bin/bash

TRACETEST="tracetest"

TEST_DEFINITION=$1
INPUT_FILE=$2

testID=$(cat "${TEST_DEFINITION}" | grep -E "^\s+id:" | tr -s ' ' | cut -d ' ' -f 3)

tmpDir=$(mktemp -d)
tmpDirName=$(basename $tmpDir)
mkdir $tmpDirName

varNames=()
line=0
while IFS=, read -ra fields; do
  # read headers
  if [ "$line" -eq "0" ]; then
    varNames=("${fields[@]}")
    line=1
    continue
  fi

  # process each line
  envFile="$tmpDirName/$line"
  environmentID="${testID}-iteration-$line"

  # create the env file
  cat <<EOF > "$envFile"
type: VariableSet
spec:
  name: ${environmentID}
  values:
EOF
  # populate env file values
  for ((i=0; i < ${#fields[@]}; i++)); do
    cat <<EOF >> "$envFile"
  - key: "${varNames[i]}"
    value: "${fields[i]}"
EOF
  done

  # run test
  echo "Running test"

  $TRACETEST apply variableset --file "${envFile}" > /dev/null
  $TRACETEST run test --file "${TEST_DEFINITION}" --vars "${envFile}"


  ((line=line+1))
done < "$INPUT_FILE"

rm -rf $tmpDirName
