NAMESPACE=$1
API_KEY=$2
AGENT_VERSION="${3:-latest}"
FILE_PATH="https://raw.githubusercontent.com/kubeshop/tracetest/k8s/agent/deploy-agent.yaml"

showUsageAndExit() {
    echo "Usage: ./script <namespace> <api-key> (<version>)?"
    echo "Examples:"
    echo "./script tracetest my-api-key"
    echo "./script my-namespace my-api-key v0.13.9"

    exit 1
}

if [ -z "${NAMESPACE}" ]; then
    echo "Error: Namespace is required"
    showUsageAndExit
fi

if [ -z "${API_KEY}" ]; then
    echo "Error: API key is required"
    showUsageAndExit
fi

kubectl create -n $NAMESPACE secret generic tracetest-agent-secret --from-literal=api-key=$API_KEY
MANIFEST=$(curl $FILE_PATH)
MANIFEST_WITH_VERSION=$(echo $MANIFEST | sed "s/:TAG/:$AGENT_VERSION/g")

echo $MANIFEST_WITH_VERSION | kubectl apply -n $NAMESPACE
