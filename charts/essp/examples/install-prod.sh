#!/bin/bash
# Example script to install ESSP in production environment

set -e

NAMESPACE="essp-prod"
RELEASE_NAME="essp"
CHART_PATH="./charts/essp"
VERSION="v1.0.0"

echo "Installing ESSP in production environment..."
echo "Namespace: $NAMESPACE"
echo "Release: $RELEASE_NAME"
echo "Version: $VERSION"
echo ""

# Prompt for confirmation
read -p "Are you sure you want to deploy to production? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
  echo "Deployment cancelled."
  exit 0
fi

# Create namespace if it doesn't exist
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Install/upgrade the chart
helm upgrade --install $RELEASE_NAME $CHART_PATH \
  --namespace $NAMESPACE \
  --values $CHART_PATH/values-prod.yaml \
  --set imsApi.image.tag=$VERSION \
  --set ssotSchool.image.tag=$VERSION \
  --set ssotDevices.image.tag=$VERSION \
  --set ssotParts.image.tag=$VERSION \
  --set syncWorker.image.tag=$VERSION \
  --wait \
  --timeout 10m \
  --atomic

echo ""
echo "Deployment complete!"
echo ""
echo "To check the status:"
echo "  kubectl get pods -n $NAMESPACE"
echo "  kubectl get hpa -n $NAMESPACE"
echo ""
echo "To view logs:"
echo "  kubectl logs -n $NAMESPACE -l app.kubernetes.io/component=ims-api -f"
