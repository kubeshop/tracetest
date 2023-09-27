NAMESPACE=$1
API_KEY=$2

kubectl create -n $NAMESPACE secret generic tracetest-agent-secret --from-literal=api-key=$API_KEY
kubectl apply -n $NAMESPACE -f https://raw.githubusercontent.com/kubeshop/tracetest/main/k8s/agent/deploy-agent.yml
