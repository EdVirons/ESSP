# ESSP Helm Chart Validation Guide

This document provides steps to validate the ESSP Helm chart before deployment.

## Pre-Deployment Validation

### 1. Chart Linting

Validate chart syntax and best practices:

```bash
helm lint charts/essp
```

Expected output: No errors or warnings.

### 2. Template Rendering

Test template rendering without installing:

```bash
# Render with default values
helm template essp charts/essp

# Render with dev values
helm template essp charts/essp -f charts/essp/values-dev.yaml

# Render with staging values
helm template essp charts/essp -f charts/essp/values-staging.yaml

# Render with production values
helm template essp charts/essp -f charts/essp/values-prod.yaml
```

### 3. Dry Run

Perform a dry-run installation:

```bash
helm install essp charts/essp \
  -f charts/essp/values-dev.yaml \
  --dry-run --debug \
  --namespace essp-dev
```

### 4. Values Schema Validation

If using Helm 3.8+, the values schema will be automatically validated:

```bash
helm install essp charts/essp \
  --set imsApi.replicaCount=-1 \
  --dry-run
```

This should fail with a validation error.

## Chart Structure Validation

### 5. Required Files Checklist

- [x] Chart.yaml - Chart metadata
- [x] values.yaml - Default values
- [x] values-dev.yaml - Development values
- [x] values-staging.yaml - Staging values
- [x] values-prod.yaml - Production values
- [x] values.schema.json - Values schema
- [x] README.md - Documentation
- [x] QUICKSTART.md - Quick start guide
- [x] CHANGELOG.md - Version history
- [x] .helmignore - Ignore patterns

### 6. Template Files Checklist

Core Templates:
- [x] templates/_helpers.tpl - Helper functions
- [x] templates/NOTES.txt - Post-install notes
- [x] templates/namespace.yaml - Namespace
- [x] templates/serviceaccount.yaml - Service account
- [x] templates/configmap.yaml - ConfigMap
- [x] templates/secrets.yaml - Secrets
- [x] templates/ingress.yaml - Ingress
- [x] templates/networkpolicy.yaml - Network policy

IMS API:
- [x] templates/ims-api/deployment.yaml
- [x] templates/ims-api/service.yaml
- [x] templates/ims-api/hpa.yaml
- [x] templates/ims-api/pdb.yaml

SSOT Services:
- [x] templates/ssot-school/deployment.yaml
- [x] templates/ssot-school/service.yaml
- [x] templates/ssot-devices/deployment.yaml
- [x] templates/ssot-devices/service.yaml
- [x] templates/ssot-parts/deployment.yaml
- [x] templates/ssot-parts/service.yaml

Workers:
- [x] templates/sync-worker/deployment.yaml

### 7. Examples Checklist

- [x] examples/install-dev.sh - Dev installation script
- [x] examples/install-prod.sh - Prod installation script
- [x] examples/values-custom.yaml - Custom values example

## Configuration Validation

### 8. Environment-Specific Values

Verify each environment has appropriate settings:

```bash
# Development - should have:
# - Lower resources
# - Debug logging
# - Single replicas
# - Autoscaling disabled
grep -A 5 "replicaCount\|logLevel\|autoscaling" charts/essp/values-dev.yaml

# Production - should have:
# - Higher resources
# - Production logging
# - Multiple replicas
# - Autoscaling enabled
grep -A 5 "replicaCount\|logLevel\|autoscaling" charts/essp/values-prod.yaml
```

### 9. Security Validation

Check security settings:

```bash
# Verify security contexts are set
grep -A 3 "securityContext" charts/essp/values.yaml

# Verify non-root user
grep "runAsNonRoot\|runAsUser" charts/essp/values.yaml

# Verify capabilities are dropped
grep -A 2 "capabilities" charts/essp/values.yaml
```

### 10. Resource Limits

Ensure all services have resource limits:

```bash
# Check IMS API resources
yq '.imsApi.resources' charts/essp/values.yaml

# Check all services have resources defined
for svc in imsApi ssotSchool ssotDevices ssotParts syncWorker; do
  echo "Checking $svc..."
  yq ".$svc.resources" charts/essp/values.yaml
done
```

## Functional Validation

### 11. Test Installation (Kind/Minikube)

Create a test cluster and install:

```bash
# Using Kind
kind create cluster --name essp-test

# Install prerequisites (PostgreSQL, Redis, NATS, MinIO)
# ... install infrastructure services ...

# Install ESSP
helm install essp charts/essp \
  -f charts/essp/values-dev.yaml \
  -n essp-dev --create-namespace

# Verify deployment
kubectl get pods -n essp-dev
kubectl get svc -n essp-dev
kubectl get hpa -n essp-dev

# Clean up
kind delete cluster --name essp-test
```

### 12. Verify All Components Deploy

Check that all enabled services are created:

```bash
# After installation
kubectl get deployments -n essp-dev
# Should show: ims-api, ssot-school, ssot-devices, ssot-parts, sync-worker

kubectl get services -n essp-dev
# Should show services for: ims-api, ssot-school, ssot-devices, ssot-parts

kubectl get hpa -n essp-dev
# Should show HPA for ims-api (if enabled)

kubectl get pdb -n essp-dev
# Should show PDB for ims-api (if enabled)
```

### 13. Health Check Validation

Verify health endpoints:

```bash
# Port forward to IMS API
kubectl port-forward -n essp-dev svc/essp-ims-api 8080:8080 &

# Test health endpoint
curl http://localhost:8080/health

# Test readiness endpoint
curl http://localhost:8080/ready

# Kill port-forward
pkill -f "port-forward.*8080:8080"
```

### 14. Configuration Validation

Check environment variables are correctly set:

```bash
# Check database configuration
kubectl exec -n essp-dev deployment/essp-ims-api -- env | grep DB_

# Check Redis configuration
kubectl exec -n essp-dev deployment/essp-ims-api -- env | grep REDIS_

# Check NATS configuration
kubectl exec -n essp-dev deployment/essp-ims-api -- env | grep NATS_

# Check MinIO configuration
kubectl exec -n essp-dev deployment/essp-ims-api -- env | grep MINIO_
```

### 15. Secrets Validation

Verify secrets are created correctly:

```bash
# Check secret exists
kubectl get secret essp-secrets -n essp-dev

# Verify secret has required keys
kubectl get secret essp-secrets -n essp-dev -o jsonpath='{.data}' | jq 'keys'

# Expected keys:
# - db-username
# - db-password
# - redis-password (if set)
# - nats-username (if set)
# - nats-password (if set)
# - minio-access-key
# - minio-secret-key
# - jwt-secret
```

## Upgrade Validation

### 16. Test Upgrade Process

```bash
# Initial install with v1.0.0
helm install essp charts/essp \
  --set imsApi.image.tag=v1.0.0 \
  -n essp-test --create-namespace

# Upgrade to v1.1.0
helm upgrade essp charts/essp \
  --set imsApi.image.tag=v1.1.0 \
  -n essp-test

# Verify upgrade
helm history essp -n essp-test

# Test rollback
helm rollback essp -n essp-test

# Verify rollback
helm history essp -n essp-test
```

## Scaling Validation

### 17. Test Autoscaling

```bash
# Deploy with autoscaling enabled
helm install essp charts/essp \
  --set imsApi.autoscaling.enabled=true \
  -n essp-test --create-namespace

# Check HPA status
kubectl get hpa -n essp-test

# Generate load to trigger scaling
# ... apply load to the service ...

# Watch scaling events
kubectl get hpa -w -n essp-test
```

## Network Validation

### 18. Test Network Policies

```bash
# Install with network policies enabled
helm install essp charts/essp \
  --set networkPolicy.enabled=true \
  -n essp-test --create-namespace

# Verify network policies exist
kubectl get networkpolicy -n essp-test

# Test connectivity between pods
kubectl run test-pod --image=curlimages/curl -n essp-test -- sleep 3600
kubectl exec test-pod -n essp-test -- curl http://essp-ims-api:8080/health
```

## Documentation Validation

### 19. Verify Documentation

Check that all documentation is present and accurate:

- [ ] README.md has installation instructions
- [ ] QUICKSTART.md provides quick start guide
- [ ] CHANGELOG.md is up to date
- [ ] NOTES.txt provides helpful post-install information
- [ ] Examples scripts are executable and work
- [ ] All configuration options are documented

### 20. Test Example Scripts

```bash
# Test dev installation script
bash charts/essp/examples/install-dev.sh

# Verify installation
kubectl get pods -n essp-dev

# Clean up
helm uninstall essp -n essp-dev
```

## Checklist Summary

- [ ] Chart lints without errors
- [ ] Templates render correctly for all environments
- [ ] Dry-run succeeds
- [ ] All required files present
- [ ] Security contexts properly configured
- [ ] Resource limits set for all services
- [ ] Test installation succeeds
- [ ] All components deploy correctly
- [ ] Health checks work
- [ ] Environment variables correctly set
- [ ] Secrets created properly
- [ ] Upgrade/rollback works
- [ ] Autoscaling functions
- [ ] Network policies work (if enabled)
- [ ] Documentation is complete
- [ ] Example scripts work

## Validation Report

Date: _______________
Validator: _______________
Chart Version: _______________
Kubernetes Version: _______________

Results:
- [ ] All checks passed
- [ ] Minor issues (list below)
- [ ] Major issues (list below)

Issues Found:
_______________________________________________________
_______________________________________________________
_______________________________________________________

Recommendations:
_______________________________________________________
_______________________________________________________
_______________________________________________________
