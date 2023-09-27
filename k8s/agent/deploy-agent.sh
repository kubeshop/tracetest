NAMESPACE=$1
API_KEY=$2
AGENT_VERSION="${3:-latest}"
FILE_PATH="https://raw.githubusercontent.com/kubeshop/tracetest/main/k8s/agent/deploy-agent.yaml"

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
curl $FILE_PATH | sed "s/:TAG/:$AGENT_VERSION/g" | kubectl apply -n $NAMESPACE -f -
