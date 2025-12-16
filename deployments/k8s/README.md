# ESSP Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the ESSP (Educational Support Services Platform) microservices platform to a Kubernetes cluster.

## Architecture Overview

The ESSP platform consists of the following services:

- **ims-api**: Main API service for the Incident Management System
- **ssot-school**: Single Source of Truth service for school data
- **ssot-devices**: Single Source of Truth service for device data
- **ssot-parts**: Single Source of Truth service for parts data
- **sync-worker**: Background worker for synchronizing SSOT data

### External Dependencies

The services require the following external infrastructure:

- **PostgreSQL**: Database (multiple databases for different services)
- **Redis/Valkey**: Caching and session storage
- **MinIO**: S3-compatible object storage for attachments
- **NATS**: Message broker for event-driven communication

## Directory Structure

```
deployments/k8s/
├── namespace.yaml                 # ESSP namespace definition
├── configmaps/                    # Configuration for all services
│   ├── ims-api-config.yaml
│   ├── ssot-config.yaml
│   └── sync-worker-config.yaml
├── secrets/                       # Secret templates (MUST BE CONFIGURED)
│   ├── database-secrets.yaml
│   ├── redis-secrets.yaml
│   ├── minio-secrets.yaml
│   └── nats-secrets.yaml
├── deployments/                   # Deployment manifests
│   ├── ims-api-deployment.yaml
│   ├── ssot-school-deployment.yaml
│   ├── ssot-devices-deployment.yaml
│   ├── ssot-parts-deployment.yaml
│   └── sync-worker-deployment.yaml
├── services/                      # Service definitions
│   ├── ims-api-service.yaml
│   ├── ssot-school-service.yaml
│   ├── ssot-devices-service.yaml
│   └── ssot-parts-service.yaml
├── ingress.yaml                   # Ingress configuration for external access
├── pdbs/                         # PodDisruptionBudgets for HA
│   ├── ims-api-pdb.yaml
│   ├── ssot-school-pdb.yaml
│   ├── ssot-devices-pdb.yaml
│   ├── ssot-parts-pdb.yaml
│   └── sync-worker-pdb.yaml
├── hpa/                          # HorizontalPodAutoscalers
│   ├── ims-api-hpa.yaml
│   ├── ssot-school-hpa.yaml
│   ├── ssot-devices-hpa.yaml
│   └── ssot-parts-hpa.yaml
└── network-policies/             # Network security policies
    ├── default-deny.yaml
    ├── allow-dns.yaml
    ├── ims-api-policy.yaml
    ├── ssot-services-policy.yaml
    └── sync-worker-policy.yaml
```

## Prerequisites

Before deploying, ensure you have:

1. **Kubernetes Cluster** (v1.24+)
   - Properly configured `kubectl` access
   - Sufficient resources for all services

2. **External Infrastructure**
   - PostgreSQL database server with created databases:
     - `ssp_ims` (for ims-api and sync-worker)
     - `ssp_school` (for ssot-school)
     - `ssp_devices` (for ssot-devices)
     - `ssp_parts` (for ssot-parts)
   - Redis/Valkey server
   - MinIO or S3-compatible storage
   - NATS message broker

3. **Kubernetes Components**
   - NGINX Ingress Controller installed
   - cert-manager (for TLS certificates) - optional but recommended
   - Metrics Server (for HPA functionality)

4. **Container Images**
   - Build and push Docker images to your registry:
     ```bash
     docker build -t your-registry/essp/ims-api:latest -f deployments/docker/ims-api.Dockerfile .
     docker build -t your-registry/essp/ssot-school:latest -f deployments/docker/ssot-school.Dockerfile .
     docker build -t your-registry/essp/ssot-devices:latest -f deployments/docker/ssot-devices.Dockerfile .
     docker build -t your-registry/essp/ssot-parts:latest -f deployments/docker/ssot-parts.Dockerfile .
     docker build -t your-registry/essp/sync-worker:latest -f deployments/docker/sync-worker.Dockerfile .
     ```

## Configuration Steps

### 1. Update Secret Values

The secret files in `secrets/` are templates. You MUST update them with your actual credentials:

```bash
# Example: Encode your database DSN
echo -n 'postgres://user:password@host:5432/ssp_ims?sslmode=require' | base64

# Example: Encode Redis password
echo -n 'your-redis-password' | base64

# Example: Encode MinIO credentials
echo -n 'your-minio-access-key' | base64
echo -n 'your-minio-secret-key' | base64
```

Edit the following files and replace `<BASE64_ENCODED_*>` placeholders:

- `secrets/database-secrets.yaml`
- `secrets/redis-secrets.yaml`
- `secrets/minio-secrets.yaml`
- `secrets/nats-secrets.yaml` (if using NATS authentication)

### 2. Update ConfigMaps

Review and update `configmaps/` files with your environment-specific values:

**Key configurations to update:**

- `configmaps/ims-api-config.yaml`:
  - `ATTACHMENTS_PUBLIC_BASE_URL`: Your MinIO/S3 public URL
  - `CORS_ALLOWED_ORIGINS`: Your frontend URL
  - `AUTH_ISSUER`: Your authentication provider URL
  - `AUTH_JWKS_URL`: Your JWKS endpoint

### 3. Update Ingress Configuration

Edit `ingress.yaml` and update:

- **Hostnames**: Replace `essp.example.com` with your actual domain
- **TLS Certificate**: Update the `secretName` or cert-manager issuer
- **Ingress Class**: Ensure it matches your ingress controller

Choose one of the two ingress configurations:
- **Option 1**: Path-based routing (single domain)
- **Option 2**: Subdomain-based routing (multiple domains)

Comment out the option you don't need.

### 4. Update Container Images

In each deployment file under `deployments/`, update the image references:

```yaml
image: your-registry/essp/ims-api:v1.0.0  # Update with your registry and tag
```

## Deployment Instructions

### Step 1: Create Namespace

```bash
kubectl apply -f namespace.yaml
```

### Step 2: Deploy Secrets

```bash
kubectl apply -f secrets/
```

Verify secrets are created:
```bash
kubectl get secrets -n essp
```

### Step 3: Deploy ConfigMaps

```bash
kubectl apply -f configmaps/
```

Verify ConfigMaps:
```bash
kubectl get configmaps -n essp
```

### Step 4: Deploy Services

```bash
kubectl apply -f services/
```

Verify services:
```bash
kubectl get services -n essp
```

### Step 5: Deploy Deployments

```bash
kubectl apply -f deployments/
```

Wait for pods to be ready:
```bash
kubectl get pods -n essp -w
```

Check deployment status:
```bash
kubectl rollout status deployment/ims-api -n essp
kubectl rollout status deployment/ssot-school -n essp
kubectl rollout status deployment/ssot-devices -n essp
kubectl rollout status deployment/ssot-parts -n essp
kubectl rollout status deployment/sync-worker -n essp
```

### Step 6: Deploy Ingress

```bash
kubectl apply -f ingress.yaml
```

Verify ingress:
```bash
kubectl get ingress -n essp
kubectl describe ingress essp-ingress -n essp
```

### Step 7: Deploy PodDisruptionBudgets

```bash
kubectl apply -f pdbs/
```

### Step 8: Deploy HorizontalPodAutoscalers

Ensure Metrics Server is running:
```bash
kubectl get apiservices | grep metrics
```

Deploy HPAs:
```bash
kubectl apply -f hpa/
```

Verify HPAs:
```bash
kubectl get hpa -n essp
```

### Step 9: Deploy Network Policies (Optional but Recommended)

**WARNING**: Network policies will restrict traffic. Only apply if your cluster supports NetworkPolicy and you've configured the policies correctly.

```bash
kubectl apply -f network-policies/
```

Verify network policies:
```bash
kubectl get networkpolicies -n essp
```

## Verification

### Check All Resources

```bash
kubectl get all -n essp
```

### Check Pod Logs

```bash
# IMS API logs
kubectl logs -n essp -l app=ims-api --tail=100 -f

# SSOT School logs
kubectl logs -n essp -l app=ssot-school --tail=100 -f

# Sync Worker logs
kubectl logs -n essp -l app=sync-worker --tail=100 -f
```

### Test Health Endpoints

```bash
# Port-forward to test locally
kubectl port-forward -n essp svc/ims-api 8080:8080

# In another terminal
curl http://localhost:8080/health
curl http://localhost:8080/readyz
```

### Test Ingress

```bash
curl https://essp.example.com/api/health
```

## Database Migrations

Before the first deployment, run database migrations:

```bash
# Create a migration job or run migrations manually
kubectl run migrate-ims --rm -it --restart=Never \
  --image=your-registry/essp/ims-api:latest \
  --env="PG_DSN=$(kubectl get secret database-secrets -n essp -o jsonpath='{.data.PG_DSN}' | base64 -d)" \
  -- /app/migrate up

# Repeat for each SSOT service
```

## Scaling

### Manual Scaling

```bash
# Scale ims-api to 5 replicas
kubectl scale deployment ims-api -n essp --replicas=5

# Scale all SSOT services
kubectl scale deployment ssot-school ssot-devices ssot-parts -n essp --replicas=3
```

### Autoscaling

HPAs are configured to automatically scale based on CPU/Memory utilization:
- **ims-api**: 2-10 replicas (target: 70% CPU)
- **ssot-school/devices/parts**: 2-6 replicas (target: 70% CPU)

Monitor autoscaling:
```bash
kubectl get hpa -n essp -w
```

## Updates and Rollbacks

### Rolling Update

```bash
# Update image
kubectl set image deployment/ims-api ims-api=your-registry/essp/ims-api:v1.1.0 -n essp

# Watch rollout
kubectl rollout status deployment/ims-api -n essp
```

### Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/ims-api -n essp

# Rollback to specific revision
kubectl rollout history deployment/ims-api -n essp
kubectl rollout undo deployment/ims-api --to-revision=2 -n essp
```

## Monitoring and Debugging

### View Events

```bash
kubectl get events -n essp --sort-by='.lastTimestamp'
```

### Describe Resources

```bash
kubectl describe pod -n essp <pod-name>
kubectl describe deployment -n essp ims-api
```

### Execute Commands in Pods

```bash
kubectl exec -it -n essp <pod-name> -- /bin/sh
```

### Check Resource Usage

```bash
kubectl top pods -n essp
kubectl top nodes
```

## Troubleshooting

### Pods Not Starting

1. Check pod status: `kubectl describe pod -n essp <pod-name>`
2. Check logs: `kubectl logs -n essp <pod-name>`
3. Verify secrets exist: `kubectl get secrets -n essp`
4. Check image pull: `kubectl get events -n essp | grep Failed`

### Database Connection Issues

1. Verify database credentials in secrets
2. Check network connectivity from pods:
   ```bash
   kubectl run -it --rm debug --image=postgres:16 --restart=Never -- psql "$DB_URL"
   ```

### Service Not Accessible

1. Check service endpoints: `kubectl get endpoints -n essp`
2. Verify pods are ready: `kubectl get pods -n essp`
3. Check ingress: `kubectl describe ingress -n essp essp-ingress`

### HPA Not Scaling

1. Verify Metrics Server is running: `kubectl get apiservices | grep metrics`
2. Check HPA status: `kubectl describe hpa -n essp ims-api-hpa`
3. Verify metrics are available: `kubectl top pods -n essp`

## Security Best Practices

1. **Secrets Management**: Use external secret managers (e.g., HashiCorp Vault, AWS Secrets Manager)
2. **RBAC**: Implement proper Role-Based Access Control
3. **Network Policies**: Keep network policies enabled in production
4. **Image Security**: Scan images for vulnerabilities before deployment
5. **TLS**: Always use TLS/HTTPS in production
6. **Resource Limits**: Keep resource limits to prevent resource exhaustion
7. **Pod Security**: Use Pod Security Standards/Policies

## Production Checklist

- [ ] All secrets properly configured with production credentials
- [ ] TLS certificates configured (cert-manager or manual)
- [ ] Resource requests/limits tuned for production workload
- [ ] HPA thresholds adjusted based on load testing
- [ ] Monitoring and alerting configured
- [ ] Log aggregation set up
- [ ] Backup strategy for databases
- [ ] Disaster recovery plan in place
- [ ] Network policies tested and enabled
- [ ] RBAC policies configured
- [ ] Regular security scanning of images

## Clean Up

To remove all ESSP resources:

```bash
# Delete all resources in namespace
kubectl delete namespace essp

# Or delete resources individually
kubectl delete -f network-policies/
kubectl delete -f hpa/
kubectl delete -f pdbs/
kubectl delete -f ingress.yaml
kubectl delete -f deployments/
kubectl delete -f services/
kubectl delete -f configmaps/
kubectl delete -f secrets/
kubectl delete -f namespace.yaml
```

## Support and Documentation

For more information about the ESSP platform, see:
- Main README: `/home/pato/opt/ESSP/README.md`
- Docker deployment: `/home/pato/opt/ESSP/deployments/docker/`
- Service documentation: `/home/pato/opt/ESSP/docs/`

## References

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [cert-manager](https://cert-manager.io/docs/)
- [Horizontal Pod Autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
