#!/bin/bash

showUsageAndExit() {
  echo "Usage: ./deploy-agent.sh <namespace> <api-key> [<version>] [--skip-verify] [--server-url <url>]"
  echo "Examples:"
  echo "  ./deploy-agent.sh tracetest my-api-key"
  echo "  ./deploy-agent.sh my-namespace my-api-key v0.13.9"
  echo "  ./deploy-agent.sh my-namespace my-api-key --skip-verify --server-url https://tracetest.my-domain.com"
  echo ""
  echo "Options:"
  echo "  --skip-verify        Skip SSL certificate verification"
  echo "  --server-url <url>   Specify the server URL"
  echo "  -h, --help           Show this help message and exit"
  echo ""
  echo "Description:"
  echo "  This script deploys the tracetest agent in a Kubernetes cluster."
  echo "  It requires a namespace and an API key to authenticate with the server."
  echo "  The agent version can be optionally specified. If not provided, the latest version will be used."
  echo "  Additional options can be used to customize the deployment."
  exit 1
}

FILE_PATH="https://raw.githubusercontent.com/kubeshop/tracetest/main/k8s/agent/deploy-agent.yaml"

NAMESPACE=""
API_KEY=""
AGENT_VERSION="latest"
SKIP_VERIFY=false
SERVER_URL=""
ENVIRONMENT=""

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-verify)
            SKIP_VERIFY=true
            shift
            ;;
        --server-url)
            SERVER_URL="$2"
            shift
            shift
            ;;
        --environment)
            ENVIRONMENT="$2"
            shift
            shift
            ;;
        -h|--help)
            showUsageAndExit
            ;;
        -*|--*)
            echo "Invalid flag: $1"
            showUsageAndExit
            ;;
        *)
            POSITIONAL_ARGS+=("$1")
            shift
            ;;
    esac
done

if [[ ${#POSITIONAL_ARGS[@]} -ge 2 ]]; then
    NAMESPACE=${POSITIONAL_ARGS[0]}
    API_KEY=${POSITIONAL_ARGS[1]}
fi

if [[ ${#POSITIONAL_ARGS[@]} -ge 3 ]]; then
    AGENT_VERSION=${POSITIONAL_ARGS[2]}
fi

if [ -z "${NAMESPACE}" ]; then
    echo "Error: Namespace is required"
    showUsageAndExit
fi

if [ -z "${API_KEY}" ]; then
    echo "Error: API key is required"
    showUsageAndExit
fi

extraCmd=()
if [ "$SKIP_VERIFY" == "true" ]; then
  extraCmd+=("--skip-verify")
fi

if [ -n "$SERVER_URL" ]; then
  extraCmd+=("--server-url" "$SERVER_URL")
fi

if [ -n "$ENVIRONMENT" ]; then
  extraCmd+=("--environment" "$ENVIRONMENT")
fi


kubectl create -n $NAMESPACE secret generic tracetest-agent-secret --from-literal=api-key=$API_KEY
curl $FILE_PATH \
  | sed "s/:TAG/:$AGENT_VERSION/g" \
  | sed "$(if [ ${#extraCmd[@]} -eq 0 ]; then echo '/EXTRA_CMD/d'; else echo "s|EXTRA_CMD|$(printf "\"%s\"," "${extraCmd[@]}" | sed 's/,$//')|g"; fi)" \
  | kubectl apply -n $NAMESPACE -f -
