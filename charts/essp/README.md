# ESSP Helm Chart

This Helm chart deploys the EdVirons School Services Platform (ESSP) microservices on a Kubernetes cluster.

## Overview

The ESSP platform consists of the following microservices:
- **ims-api**: Main API gateway service
- **ssot-school**: School source of truth service
- **ssot-devices**: Devices source of truth service
- **ssot-parts**: Parts source of truth service
- **sync-worker**: Background synchronization worker

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- PostgreSQL database
- Redis instance
- NATS messaging server
- MinIO object storage

## Installation

### Basic Installation

```bash
# Install with default values
helm install essp ./charts/essp

# Install with custom release name
helm install my-essp ./charts/essp
```

### Environment-Specific Installation

#### Development Environment

```bash
helm install essp ./charts/essp \
  -f charts/essp/values-dev.yaml \
  --namespace essp-dev \
  --create-namespace
```

#### Staging Environment

```bash
helm install essp ./charts/essp \
  -f charts/essp/values-staging.yaml \
  --namespace essp-staging \
  --create-namespace
```

#### Production Environment

```bash
helm install essp ./charts/essp \
  -f charts/essp/values-prod.yaml \
  --namespace essp-prod \
  --create-namespace
```

### Installation with Custom Values

```bash
helm install essp ./charts/essp \
  --set imsApi.image.tag=v1.0.0 \
  --set global.imageRegistry=myregistry.io/
```

## Configuration

### Global Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.namespace` | Kubernetes namespace | `essp` |
| `global.imageRegistry` | Container image registry prefix | `""` |
| `global.imagePullSecrets` | Image pull secrets | `[]` |

### IMS API Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `imsApi.enabled` | Enable IMS API service | `true` |
| `imsApi.replicaCount` | Number of replicas | `2` |
| `imsApi.image.repository` | Image repository | `essp/ims-api` |
| `imsApi.image.tag` | Image tag | `latest` |
| `imsApi.resources.requests.cpu` | CPU request | `100m` |
| `imsApi.resources.requests.memory` | Memory request | `128Mi` |
| `imsApi.autoscaling.enabled` | Enable HPA | `true` |
| `imsApi.autoscaling.minReplicas` | Minimum replicas | `2` |
| `imsApi.autoscaling.maxReplicas` | Maximum replicas | `10` |

### Database Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `database.host` | PostgreSQL host | `postgresql` |
| `database.port` | PostgreSQL port | `5432` |
| `database.name` | Database name | `ssp_ims` |
| `database.sslMode` | SSL mode | `disable` |
| `secrets.database.username` | Database username | `""` |
| `secrets.database.password` | Database password | `""` |

### Redis Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `redis.host` | Redis host | `redis` |
| `redis.port` | Redis port | `6379` |
| `redis.db` | Redis database | `0` |

### NATS Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nats.url` | NATS URL | `nats://nats:4222` |
| `nats.cluster` | NATS cluster name | `essp-cluster` |

### MinIO Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `minio.endpoint` | MinIO endpoint | `minio:9000` |
| `minio.bucket` | Bucket name | `essp-attachments` |
| `minio.useSSL` | Use SSL | `false` |

### Ingress Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable ingress | `true` |
| `ingress.className` | Ingress class | `nginx` |
| `ingress.hosts[0].host` | Hostname | `api.essp.local` |

## Secrets Management

### Using Values File

Create a custom values file with secrets:

```yaml
# secrets.yaml
secrets:
  database:
    username: "dbuser"
    password: "dbpassword"
  redis:
    password: "redispassword"
  minio:
    accessKey: "minioaccess"
    secretKey: "miniosecret"
  jwt:
    secret: "jwtsecret"
```

Install with secrets:

```bash
helm install essp ./charts/essp -f secrets.yaml
```

### Using External Secrets Operator

For production environments, use External Secrets Operator:

```yaml
# Disable built-in secrets
secrets:
  database:
    username: ""
    password: ""

# Create ExternalSecret resources separately
```

### Using kubectl create secret

```bash
kubectl create secret generic essp-secrets \
  --from-literal=db-username=myuser \
  --from-literal=db-password=mypassword \
  --namespace essp
```

## Upgrade

### Standard Upgrade

```bash
helm upgrade essp ./charts/essp
```

### Upgrade with New Values

```bash
helm upgrade essp ./charts/essp \
  -f charts/essp/values-prod.yaml \
  --set imsApi.image.tag=v1.1.0
```

### Rollback

```bash
# List release history
helm history essp

# Rollback to previous revision
helm rollback essp

# Rollback to specific revision
helm rollback essp 2
```

## Uninstallation

```bash
# Uninstall release
helm uninstall essp

# Uninstall and delete namespace
helm uninstall essp -n essp
kubectl delete namespace essp
```

## Monitoring and Observability

### Health Checks

All services expose health endpoints:
- Liveness: `/health`
- Readiness: `/ready`

### Logs

```bash
# View logs for IMS API
kubectl logs -l app.kubernetes.io/component=ims-api -n essp

# View logs for all services
kubectl logs -l app.kubernetes.io/name=essp -n essp --all-containers

# Follow logs
kubectl logs -f -l app.kubernetes.io/component=ims-api -n essp
```

### Metrics

If Prometheus is installed, metrics are automatically scraped from service endpoints.

## Scaling

### Manual Scaling

```bash
# Scale IMS API to 5 replicas
kubectl scale deployment essp-ims-api --replicas=5 -n essp
```

### Horizontal Pod Autoscaler

HPA is enabled by default for the IMS API service and automatically scales based on CPU and memory utilization.

```bash
# Check HPA status
kubectl get hpa -n essp

# Describe HPA
kubectl describe hpa essp-ims-api -n essp
```

## Network Policies

Network policies can be enabled to restrict traffic between pods:

```yaml
networkPolicy:
  enabled: true
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n essp
kubectl describe pod <pod-name> -n essp
```

### Check Service Endpoints

```bash
kubectl get svc -n essp
kubectl get endpoints -n essp
```

### Check Ingress

```bash
kubectl get ingress -n essp
kubectl describe ingress essp -n essp
```

### Debug Container

```bash
# Execute shell in pod
kubectl exec -it <pod-name> -n essp -- /bin/sh

# Check environment variables
kubectl exec <pod-name> -n essp -- env
```

### Common Issues

#### Pods Not Starting
- Check image pull secrets: `kubectl describe pod <pod-name> -n essp`
- Verify resource limits: `kubectl describe pod <pod-name> -n essp`
- Check logs: `kubectl logs <pod-name> -n essp`

#### Database Connection Errors
- Verify database credentials in secrets
- Check database host and port configuration
- Ensure database is accessible from cluster

#### Service Not Accessible
- Verify service is running: `kubectl get svc -n essp`
- Check ingress configuration: `kubectl describe ingress -n essp`
- Verify DNS resolution

## Development

### Validate Chart

```bash
# Lint the chart
helm lint ./charts/essp

# Dry run installation
helm install essp ./charts/essp --dry-run --debug

# Template rendering
helm template essp ./charts/essp
```

### Testing

```bash
# Install in test mode
helm install essp ./charts/essp --dry-run --debug

# Test with different values
helm install essp ./charts/essp \
  -f charts/essp/values-dev.yaml \
  --dry-run --debug
```

## Best Practices

1. **Always use specific image tags** in production instead of `latest`
2. **Store secrets securely** using external secret management systems
3. **Enable network policies** in production for better security
4. **Configure resource limits** appropriate for your workload
5. **Use separate namespaces** for different environments
6. **Enable pod disruption budgets** for critical services
7. **Monitor HPA metrics** to optimize autoscaling settings
8. **Regular backups** of database and persistent data

## Support

For issues and questions:
- GitHub Issues: [repository URL]
- Documentation: [docs URL]
- Email: dev@edvirons.com

## License

Copyright EdVirons. All rights reserved.
