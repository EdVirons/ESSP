# ESSP Helm Chart - Quick Start Guide

This guide will help you get ESSP up and running quickly in different environments.

## Prerequisites

Before you begin, ensure you have:
- Kubernetes cluster (1.19+)
- Helm 3.0+ installed
- kubectl configured to access your cluster
- Required infrastructure services (PostgreSQL, Redis, NATS, MinIO)

## Quick Installation

### Development Environment

For local development or testing:

```bash
# 1. Create namespace
kubectl create namespace essp-dev

# 2. Install with development values
helm install essp ./charts/essp \
  -f ./charts/essp/values-dev.yaml \
  -n essp-dev

# 3. Wait for pods to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=essp -n essp-dev --timeout=300s

# 4. Port forward to access API
kubectl port-forward -n essp-dev svc/essp-ims-api 8080:8080

# Access API at http://localhost:8080
```

Or use the provided script:

```bash
./charts/essp/examples/install-dev.sh
```

### Staging Environment

For staging/pre-production:

```bash
helm install essp ./charts/essp \
  -f ./charts/essp/values-staging.yaml \
  -n essp-staging \
  --create-namespace
```

### Production Environment

For production deployment:

```bash
# Create production secrets file (DO NOT commit to git)
cat > prod-secrets.yaml <<EOF
secrets:
  database:
    username: "your_db_user"
    password: "your_secure_db_password"
  redis:
    password: "your_secure_redis_password"
  minio:
    accessKey: "your_minio_access_key"
    secretKey: "your_minio_secret_key"
  jwt:
    secret: "your_jwt_secret_key"
EOF

# Install with production values and secrets
helm install essp ./charts/essp \
  -f ./charts/essp/values-prod.yaml \
  -f prod-secrets.yaml \
  -n essp-prod \
  --create-namespace

# Clean up secrets file
rm prod-secrets.yaml
```

Or use the provided script:

```bash
./charts/essp/examples/install-prod.sh
```

## Verify Installation

```bash
# Check pod status
kubectl get pods -n essp

# Check services
kubectl get svc -n essp

# Check ingress
kubectl get ingress -n essp

# View logs
kubectl logs -l app.kubernetes.io/component=ims-api -n essp
```

## Common Customizations

### Custom Image Tags

```bash
helm install essp ./charts/essp \
  --set imsApi.image.tag=v1.2.3 \
  --set ssotSchool.image.tag=v1.2.3
```

### Custom Database

```bash
helm install essp ./charts/essp \
  --set database.host=my-postgres.example.com \
  --set database.port=5432 \
  --set database.name=my_database
```

### Custom Ingress Host

```bash
helm install essp ./charts/essp \
  --set ingress.hosts[0].host=api.mycompany.com
```

### Disable Specific Services

```bash
helm install essp ./charts/essp \
  --set ssotParts.enabled=false \
  --set syncWorker.enabled=false
```

## Accessing the Application

### Via Ingress (Production)

If ingress is configured with a hostname:

```bash
curl https://api.essp.example.com/health
```

### Via Port Forwarding (Development)

```bash
# Forward IMS API
kubectl port-forward -n essp svc/essp-ims-api 8080:8080

# Access API
curl http://localhost:8080/health
```

### Via LoadBalancer (Cloud)

```bash
# Change service type to LoadBalancer
helm upgrade essp ./charts/essp \
  --set imsApi.service.type=LoadBalancer

# Get external IP
kubectl get svc essp-ims-api -n essp
```

## Upgrade

```bash
# Upgrade to new version
helm upgrade essp ./charts/essp \
  --set imsApi.image.tag=v1.3.0 \
  -n essp

# Upgrade with new values file
helm upgrade essp ./charts/essp \
  -f ./charts/essp/values-prod.yaml \
  -n essp
```

## Rollback

```bash
# View release history
helm history essp -n essp

# Rollback to previous version
helm rollback essp -n essp

# Rollback to specific revision
helm rollback essp 2 -n essp
```

## Uninstall

```bash
# Uninstall the release
helm uninstall essp -n essp

# Delete namespace (optional)
kubectl delete namespace essp
```

## Troubleshooting

### Pods Not Starting

```bash
# Check pod details
kubectl describe pod <pod-name> -n essp

# Check events
kubectl get events -n essp --sort-by='.lastTimestamp'
```

### Service Not Accessible

```bash
# Check service endpoints
kubectl get endpoints -n essp

# Test service from within cluster
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
  curl http://essp-ims-api.essp.svc.cluster.local:8080/health
```

### Database Connection Issues

```bash
# Check secrets
kubectl get secret essp-secrets -n essp -o yaml

# Check environment variables
kubectl exec <pod-name> -n essp -- env | grep DB_
```

## Next Steps

- Review the full [README.md](./README.md) for detailed configuration options
- Check [CHANGELOG.md](./CHANGELOG.md) for version history
- See [examples/values-custom.yaml](./examples/values-custom.yaml) for advanced configuration
- Configure monitoring and alerting
- Set up CI/CD pipelines for automated deployments
- Implement backup and disaster recovery procedures

## Getting Help

- Documentation: See [README.md](./README.md)
- Issues: Report bugs and request features via GitHub Issues
- Email: dev@edvirons.com

## Security Notes

- Always use secrets management in production (e.g., sealed-secrets, external-secrets)
- Never commit secrets to version control
- Enable TLS for all external communications
- Regularly update container images
- Enable network policies in production
- Use RBAC to restrict access
