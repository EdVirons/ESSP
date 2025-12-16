#!/bin/bash
# Example script to install ESSP in development environment

set -e

NAMESPACE="essp-dev"
RELEASE_NAME="essp"
CHART_PATH="./charts/essp"

echo "Installing ESSP in development environment..."
echo "Namespace: $NAMESPACE"
echo "Release: $RELEASE_NAME"
echo ""

# Create namespace if it doesn't exist
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Install the chart
helm upgrade --install $RELEASE_NAME $CHART_PATH \
  --namespace $NAMESPACE \
  --values $CHART_PATH/values-dev.yaml \
  --set imsApi.image.tag=dev \
  --set ssotSchool.image.tag=dev \
  --set ssotDevices.image.tag=dev \
  --set ssotParts.image.tag=dev \
  --set syncWorker.image.tag=dev \
  --wait \
  --timeout 5m

echo ""
echo "Installation complete!"
echo ""
echo "To check the status:"
echo "  kubectl get pods -n $NAMESPACE"
echo ""
echo "To access the API:"
echo "  kubectl port-forward -n $NAMESPACE svc/$RELEASE_NAME-ims-api 8080:8080"
echo ""
echo "To view logs:"
echo "  kubectl logs -n $NAMESPACE -l app.kubernetes.io/component=ims-api -f"
