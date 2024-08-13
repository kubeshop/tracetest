#!/bin/sh

set -e

TOKEN=$TRACETEST_TOKEN
ENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID
ACCESS_ID=$SUMO_LOGIC_ACCESS_ID
ACCESS_KEY=$SUMO_LOGIC_ACCESS_KEY

apply() {
  echo "Configuring TraceTest"
  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID

  echo "
---
type: DataStore
spec:
  id: current
  name: Sumo Logic
  type: sumologic
  sumologic:
    # The URL will differ based on your location. View this
    # docs page to figure out which URL you need:
    # https://help.sumologic.com/docs/api/getting-started/#which-endpoint-should-i-should-use
    url: "https://api.sumologic.com/api/"
    # Create your ID and Key under Administration > Security > Access Keys
    # in your Sumo Logic account:
    # https://help.sumologic.com/docs/manage/security/access-keys/#create-your-access-key
    accessID: ${ACCESS_ID}
    accessKey: ${ACCESS_KEY}
" > /resources/datastore.yaml

  echo "Applying Resources"
  tracetest apply datastore -f /resources/datastore.yaml
  tracetest apply test -f /resources/test.yaml
}

apply
