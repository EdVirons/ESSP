# ESSP Helm Chart - Getting Started

Welcome to the ESSP Helm Chart! This guide will help you get started quickly.

## Quick Links

- **Full Documentation**: [README.md](./README.md)
- **Quick Start Guide**: [QUICKSTART.md](./QUICKSTART.md)
- **Validation Guide**: [VALIDATION.md](./VALIDATION.md)
- **Change Log**: [CHANGELOG.md](./CHANGELOG.md)

## Installation in 3 Steps

### Step 1: Prepare Your Environment

Ensure you have:
- Kubernetes cluster running (1.19+)
- Helm 3.0+ installed
- kubectl configured
- Infrastructure services (PostgreSQL, Redis, NATS, MinIO)

### Step 2: Choose Your Environment

**Development** (recommended for first time):
```bash
helm install essp ./charts/essp \
  -f ./charts/essp/values-dev.yaml \
  -n essp-dev --create-namespace
```

**Staging**:
```bash
helm install essp ./charts/essp \
  -f ./charts/essp/values-staging.yaml \
  -n essp-staging --create-namespace
```

**Production**:
```bash
# Create secrets file first (don't commit to git!)
helm install essp ./charts/essp \
  -f ./charts/essp/values-prod.yaml \
  -f prod-secrets.yaml \
  -n essp-prod --create-namespace
```

### Step 3: Verify Installation

```bash
# Check pods
kubectl get pods -n essp-dev

# Check services
kubectl get svc -n essp-dev

# View logs
kubectl logs -l app.kubernetes.io/component=ims-api -n essp-dev
```

## Access Your Application

### Port Forward (Development)
```bash
kubectl port-forward -n essp-dev svc/essp-ims-api 8080:8080
curl http://localhost:8080/health
```

### Via Ingress (Production)
```bash
curl https://api.essp.example.com/health
```

## Common Commands

### View Release Info
```bash
helm list -n essp-dev
helm status essp -n essp-dev
```

### Upgrade
```bash
helm upgrade essp ./charts/essp \
  -f ./charts/essp/values-dev.yaml \
  -n essp-dev
```

### Rollback
```bash
helm rollback essp -n essp-dev
```

### Uninstall
```bash
helm uninstall essp -n essp-dev
```

## Customization

### Override Values
```bash
helm install essp ./charts/essp \
  --set imsApi.image.tag=v1.2.3 \
  --set imsApi.replicaCount=3
```

### Use Custom Values File
```bash
helm install essp ./charts/essp \
  -f ./charts/essp/values-prod.yaml \
  -f my-custom-values.yaml
```

See [examples/values-custom.yaml](./examples/values-custom.yaml) for examples.

## Troubleshooting

### Pods Not Starting?
```bash
kubectl describe pod <pod-name> -n essp-dev
kubectl logs <pod-name> -n essp-dev
```

### Service Not Accessible?
```bash
kubectl get svc -n essp-dev
kubectl get endpoints -n essp-dev
```

### Configuration Issues?
```bash
kubectl get configmap essp-config -n essp-dev -o yaml
kubectl get secret essp-secrets -n essp-dev -o yaml
```

## What's Included

This Helm chart deploys:
- **IMS API**: Main API gateway with autoscaling
- **SSOT School**: School data service
- **SSOT Devices**: Device management service
- **SSOT Parts**: Parts inventory service
- **Sync Worker**: Background synchronization worker

## Key Features

- Multi-environment support (dev/staging/prod)
- High availability with autoscaling
- Security hardened (non-root, read-only filesystem)
- Health checks and monitoring
- Network policies support
- Comprehensive documentation

## Environment Differences

| Feature | Development | Staging | Production |
|---------|-------------|---------|------------|
| Replicas | 1 | 2 | 3 |
| Autoscaling | Disabled | Enabled | Enabled |
| Resources | Low | Medium | High |
| Logging | Debug | Info | Warn |
| TLS | No | Yes | Yes |

## Next Steps

1. Read the [QUICKSTART.md](./QUICKSTART.md) for detailed installation
2. Review [README.md](./README.md) for all configuration options
3. Check [VALIDATION.md](./VALIDATION.md) before production deployment
4. Customize values files for your environment
5. Set up proper secrets management
6. Configure monitoring and alerting

## Need Help?

- Documentation: See [README.md](./README.md)
- Examples: Check [examples/](./examples/) directory
- Validation: Review [VALIDATION.md](./VALIDATION.md)
- Contact: dev@edvirons.com

## Quick Reference

### File Structure
```
charts/essp/
├── Chart.yaml              # Chart metadata
├── values*.yaml            # Configuration files
├── templates/              # Kubernetes manifests
├── examples/               # Example scripts
└── *.md                    # Documentation
```

### Important Values

**Image Configuration**:
```yaml
imsApi:
  image:
    repository: essp/ims-api
    tag: latest
```

**Database**:
```yaml
database:
  host: postgresql
  port: 5432
  name: ssp_ims
```

**Resources**:
```yaml
imsApi:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
```

### Useful Scripts

Development installation:
```bash
./examples/install-dev.sh
```

Production installation:
```bash
./examples/install-prod.sh
```

## Version Information

- **Chart Version**: 0.1.0
- **App Version**: 1.0.0
- **Kubernetes**: 1.19+
- **Helm**: 3.0+

---

**Happy Deploying!** 🚀
