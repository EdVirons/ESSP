# ESSP Platform Deployment Runbook

Version: 1.0.0
Last Updated: 2025-12-12
Platform Version: v1.0.0

## Table of Contents

1. [Overview](#1-overview)
2. [Prerequisites](#2-prerequisites)
3. [Infrastructure Setup](#3-infrastructure-setup)
4. [Secret Management](#4-secret-management)
5. [Deployment Procedures](#5-deployment-procedures)
6. [Database Migrations](#6-database-migrations)
7. [Configuration](#7-configuration)
8. [Health Checks](#8-health-checks)
9. [Scaling](#9-scaling)
10. [Monitoring & Alerting](#10-monitoring--alerting)
11. [Troubleshooting](#11-troubleshooting)
12. [Rollback Procedures](#12-rollback-procedures)
13. [Disaster Recovery](#13-disaster-recovery)
14. [Maintenance](#14-maintenance)

---

## 1. Overview

### 1.1 Architecture Diagram

```
                                    ┌─────────────────┐
                                    │   Internet      │
                                    └────────┬────────┘
                                             │
                                    ┌────────▼────────┐
                                    │  Ingress/TLS    │
                                    │  (NGINX)        │
                                    └────────┬────────┘
                                             │
                        ┌────────────────────┼────────────────────┐
                        │                    │                    │
                ┌───────▼──────┐    ┌───────▼──────┐    ┌───────▼──────┐
                │   IMS API    │    │ SSOT Services│    │ Sync Worker  │
                │  (2-20 pods) │    │ (School,Dev, │    │  (1-2 pods)  │
                │              │    │  Parts)      │    │              │
                └───────┬──────┘    └───────┬──────┘    └───────┬──────┘
                        │                   │                    │
                        └───────────────────┼────────────────────┘
                                           │
                        ┌──────────────────┼──────────────────┐
                        │                  │                  │
                ┌───────▼──────┐  ┌───────▼──────┐  ┌───────▼──────┐
                │  PostgreSQL  │  │  Redis/Valkey│  │     NATS     │
                │  (Primary +  │  │  (Cluster)   │  │  (Cluster)   │
                │   Replicas)  │  └──────────────┘  └──────────────┘
                └──────────────┘           │
                        │          ┌───────▼──────┐
                        │          │    MinIO     │
                        │          │ (S3-compat)  │
                        │          └──────────────┘
                        │
                ┌───────▼──────┐
                │  Keycloak    │
                │  (Auth/OIDC) │
                └──────────────┘

                    ┌──────────────────────────────┐
                    │   Monitoring Stack           │
                    │  Prometheus | Grafana        │
                    │  AlertManager | Loki         │
                    └──────────────────────────────┘
```

### 1.2 Service Dependencies

| Service | Depends On | Purpose |
|---------|------------|---------|
| **ims-api** | PostgreSQL, Redis, NATS, MinIO, Keycloak, SSOT Services | Main API for incident management |
| **ssot-school** | PostgreSQL | Single source of truth for school data |
| **ssot-devices** | PostgreSQL | Single source of truth for device data |
| **ssot-parts** | PostgreSQL | Single source of truth for parts data |
| **sync-worker** | PostgreSQL, NATS, SSOT Services | Background worker for data synchronization |

### 1.3 Environment Matrix

| Aspect | Development | Staging | Production |
|--------|-------------|---------|------------|
| **Namespace** | essp-dev | essp-staging | essp-prod |
| **Replicas (IMS API)** | 1 | 2 | 3-20 (HPA) |
| **Replicas (SSOT)** | 1 | 1 | 2 |
| **Replicas (Sync)** | 1 | 1 | 2 |
| **HPA** | Disabled | Enabled | Enabled |
| **PDB** | Disabled | Enabled | Enabled |
| **Network Policies** | Disabled | Enabled | Enabled |
| **TLS/SSL** | Optional | Required | Required |
| **Auth** | Disabled | Enabled | Enabled |
| **Log Level** | debug | info | warn |
| **Resource Requests** | 25-50m CPU | 50-100m CPU | 100-200m CPU |
| **Resource Limits** | 100-200m CPU | 200-500m CPU | 500-1000m CPU |
| **Database SSL** | Disabled | Optional | Required |
| **Image Pull Policy** | Always | IfNotPresent | IfNotPresent |
| **Sync Interval** | 1m | 5m | 5m |

---

## 2. Prerequisites

### 2.1 Kubernetes Cluster Requirements

#### Minimum Cluster Specifications

- **Kubernetes Version**: 1.24+
- **Nodes**: Minimum 3 worker nodes for production
- **CPU**: 8 cores total (production), 4 cores (staging), 2 cores (dev)
- **Memory**: 16GB total (production), 8GB (staging), 4GB (dev)
- **Storage**: 100GB SSD storage for persistent volumes
- **Network**: CNI plugin supporting NetworkPolicy (Calico, Cilium, etc.)

#### Required Add-ons

```bash
# Verify Kubernetes version
kubectl version --short

# Check cluster nodes
kubectl get nodes -o wide

# Verify metrics server (required for HPA)
kubectl get apiservices | grep metrics.k8s.io
# If not installed:
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Verify NGINX Ingress Controller
kubectl get pods -n ingress-nginx
# If not installed:
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml

# Optional: Install cert-manager for automatic TLS
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.2/cert-manager.yaml
```

### 2.2 Required Tools

Install the following tools on your deployment machine:

| Tool | Version | Purpose | Installation |
|------|---------|---------|--------------|
| **kubectl** | 1.24+ | Kubernetes CLI | [docs](https://kubernetes.io/docs/tasks/tools/) |
| **helm** | 3.12+ | Package manager | [docs](https://helm.sh/docs/intro/install/) |
| **docker** | 20.10+ | Container runtime | [docs](https://docs.docker.com/engine/install/) |
| **psql** | 14+ | PostgreSQL client | `apt install postgresql-client` |
| **redis-cli** | 7.0+ | Redis client | `apt install redis-tools` |
| **jq** | 1.6+ | JSON processor | `apt install jq` |
| **yq** | 4.0+ | YAML processor | [docs](https://github.com/mikefarah/yq) |
| **base64** | - | Encoding tool | Built-in on Linux/macOS |

#### Verify Tool Installation

```bash
# Verify kubectl
kubectl version --client

# Verify Helm
helm version

# Verify Docker
docker --version

# Verify PostgreSQL client
psql --version

# Verify Redis CLI
redis-cli --version

# Verify jq
jq --version

# Verify yq
yq --version
```

### 2.3 Access Requirements

#### Kubernetes Access

```bash
# Test cluster access
kubectl auth can-i create deployments --namespace=essp-prod
kubectl auth can-i create secrets --namespace=essp-prod

# Recommended RBAC permissions
# - deployments, services, configmaps, secrets: read, write, delete
# - pods, pods/log: read, list
# - ingress: read, write
# - horizontalpodautoscalers, poddisruptionbudgets: read, write
```

#### Container Registry Access

```bash
# Login to container registry
docker login <your-registry>

# Create pull secret for Kubernetes
kubectl create secret docker-registry regcred \
  --docker-server=<your-registry> \
  --docker-username=<username> \
  --docker-password=<password> \
  --docker-email=<email> \
  --namespace=essp-prod
```

#### Cloud Provider Access

- AWS: IAM roles for EKS, RDS, S3
- GCP: Service accounts for GKE, Cloud SQL, GCS
- Azure: Service principals for AKS, Azure Database, Blob Storage

---

## 3. Infrastructure Setup

### 3.1 PostgreSQL Deployment

#### Option 1: Managed Database (Recommended for Production)

**AWS RDS**:
```bash
# Create RDS PostgreSQL instance
aws rds create-db-instance \
  --db-instance-identifier essp-prod-db \
  --db-instance-class db.r6g.xlarge \
  --engine postgres \
  --engine-version 16.1 \
  --master-username essp_admin \
  --master-user-password <strong-password> \
  --allocated-storage 100 \
  --storage-type gp3 \
  --storage-encrypted \
  --backup-retention-period 30 \
  --multi-az \
  --vpc-security-group-ids <sg-id>

# Create databases
psql -h essp-prod-db.xxxxxx.region.rds.amazonaws.com -U essp_admin -d postgres <<EOF
CREATE DATABASE ssp_ims_prod;
CREATE DATABASE ssp_school_prod;
CREATE DATABASE ssp_devices_prod;
CREATE DATABASE ssp_parts_prod;
EOF
```

**GCP Cloud SQL**:
```bash
# Create Cloud SQL instance
gcloud sql instances create essp-prod-db \
  --database-version=POSTGRES_16 \
  --tier=db-custom-4-16384 \
  --region=us-central1 \
  --storage-type=SSD \
  --storage-size=100GB \
  --backup \
  --enable-bin-log

# Create databases
gcloud sql databases create ssp_ims_prod --instance=essp-prod-db
gcloud sql databases create ssp_school_prod --instance=essp-prod-db
gcloud sql databases create ssp_devices_prod --instance=essp-prod-db
gcloud sql databases create ssp_parts_prod --instance=essp-prod-db
```

#### Option 2: Self-Managed on Kubernetes

```bash
# Add Bitnami Helm repo
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install PostgreSQL
helm install postgresql bitnami/postgresql \
  --namespace essp-prod \
  --create-namespace \
  --set auth.username=essp_admin \
  --set auth.password=<strong-password> \
  --set auth.database=ssp_ims_prod \
  --set primary.persistence.size=100Gi \
  --set primary.resources.requests.cpu=1000m \
  --set primary.resources.requests.memory=2Gi \
  --set metrics.enabled=true

# Create additional databases
kubectl run psql-client --rm -it --restart=Never \
  --namespace essp-prod \
  --image=postgres:16 \
  --env="PGPASSWORD=<password>" \
  -- psql -h postgresql -U essp_admin -d postgres <<EOF
CREATE DATABASE ssp_school_prod;
CREATE DATABASE ssp_devices_prod;
CREATE DATABASE ssp_parts_prod;
EOF
```

#### Database Connection String Format

```
# PostgreSQL DSN format
postgres://username:password@host:5432/database?sslmode=require

# Examples:
# Production: postgres://essp_admin:password@essp-prod-db.xxxxxx.rds.amazonaws.com:5432/ssp_ims_prod?sslmode=require
# Staging: postgres://essp_admin:password@postgresql.essp-staging.svc.cluster.local:5432/ssp_ims_staging?sslmode=disable
# Dev: postgres://essp_admin:password@postgresql.essp-dev.svc.cluster.local:5432/ssp_ims_dev?sslmode=disable
```

### 3.2 Redis/Valkey Cluster

#### Using Helm (Recommended)

```bash
# Install Redis with Bitnami chart
helm install redis bitnami/redis \
  --namespace essp-prod \
  --set architecture=replication \
  --set auth.enabled=true \
  --set auth.password=<strong-redis-password> \
  --set master.persistence.size=8Gi \
  --set replica.replicaCount=2 \
  --set replica.persistence.size=8Gi \
  --set metrics.enabled=true \
  --set metrics.serviceMonitor.enabled=true

# Verify installation
kubectl get pods -n essp-prod -l app.kubernetes.io/name=redis

# Test connection
kubectl run redis-client --rm -it --restart=Never \
  --namespace essp-prod \
  --image=redis:7.2 \
  -- redis-cli -h redis-master -a <password> ping
```

#### Valkey Alternative

```bash
# Deploy Valkey (Redis fork)
helm install valkey oci://registry-1.docker.io/bitnamicharts/valkey \
  --namespace essp-prod \
  --set auth.enabled=true \
  --set auth.password=<strong-password> \
  --set master.persistence.size=8Gi
```

### 3.3 NATS Cluster

```bash
# Add NATS Helm repo
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update

# Install NATS with JetStream
helm install nats nats/nats \
  --namespace essp-prod \
  --set nats.jetstream.enabled=true \
  --set nats.jetstream.fileStore.pvc.size=10Gi \
  --set cluster.enabled=true \
  --set cluster.replicas=3 \
  --set auth.enabled=true \
  --set auth.timeout=5s

# Verify installation
kubectl get pods -n essp-prod -l app.kubernetes.io/name=nats

# Test connection
kubectl run nats-client --rm -it --restart=Never \
  --namespace essp-prod \
  --image=natsio/nats-box:latest \
  -- nats-server --version
```

#### NATS Configuration

```yaml
# Connection URL format
nats://nats.essp-prod.svc.cluster.local:4222

# With authentication
nats://username:password@nats.essp-prod.svc.cluster.local:4222
```

### 3.4 MinIO Setup

#### Using Helm

```bash
# Add MinIO Helm repo
helm repo add minio https://charts.min.io/
helm repo update

# Install MinIO
helm install minio minio/minio \
  --namespace essp-prod \
  --set mode=distributed \
  --set replicas=4 \
  --set persistence.size=50Gi \
  --set resources.requests.memory=1Gi \
  --set rootUser=minioadmin \
  --set rootPassword=<strong-minio-password> \
  --set buckets[0].name=essp-attachments-prod \
  --set buckets[0].policy=none \
  --set buckets[0].purge=false

# Verify installation
kubectl get pods -n essp-prod -l app=minio

# Access MinIO Console (optional)
kubectl port-forward -n essp-prod svc/minio-console 9001:9001
# Open http://localhost:9001
```

#### Create Buckets and Access Keys

```bash
# Port-forward to MinIO
kubectl port-forward -n essp-prod svc/minio 9000:9000

# Using MinIO Client (mc)
mc alias set essp http://localhost:9000 minioadmin <password>
mc mb essp/essp-attachments-prod
mc mb essp/essp-attachments-staging
mc mb essp/essp-attachments-dev

# Create service account for ESSP
mc admin user add essp essp_service_account <strong-password>
mc admin policy attach essp readwrite --user essp_service_account
```

#### S3-Compatible Alternative (AWS, GCS)

```bash
# AWS S3
aws s3 mb s3://essp-attachments-prod --region us-east-1

# GCS with S3 compatibility
gsutil mb -c STANDARD -l us-central1 gs://essp-attachments-prod
```

### 3.5 Keycloak Configuration

#### Install Keycloak

```bash
# Add Bitnami Helm repo
helm repo add bitnami https://charts.bitnami.com/bitnami

# Install Keycloak
helm install keycloak bitnami/keycloak \
  --namespace essp-prod \
  --set auth.adminUser=admin \
  --set auth.adminPassword=<strong-admin-password> \
  --set postgresql.enabled=true \
  --set postgresql.auth.password=<strong-db-password> \
  --set service.type=ClusterIP \
  --set ingress.enabled=true \
  --set ingress.hostname=auth.essp.example.com

# Verify installation
kubectl get pods -n essp-prod -l app.kubernetes.io/name=keycloak
```

#### Configure Keycloak for ESSP

1. **Access Keycloak Admin Console**:
   ```bash
   kubectl port-forward -n essp-prod svc/keycloak 8080:80
   # Open http://localhost:8080
   ```

2. **Create Realm**:
   - Navigate to "Add realm"
   - Name: `essp`
   - Enabled: Yes

3. **Create Client**:
   - Client ID: `essp-api`
   - Client Protocol: `openid-connect`
   - Access Type: `confidential`
   - Valid Redirect URIs: `https://api.essp.example.com/*`
   - Web Origins: `https://app.essp.example.com`

4. **Create Roles**:
   - `essp-admin`
   - `essp-technician`
   - `essp-viewer`

5. **Get JWKS URL**:
   ```
   https://auth.essp.example.com/realms/essp/protocol/openid-connect/certs
   ```

6. **Update ESSP Configuration**:
   ```yaml
   AUTH_ENABLED: "true"
   AUTH_ISSUER: "https://auth.essp.example.com/realms/essp"
   AUTH_JWKS_URL: "https://auth.essp.example.com/realms/essp/protocol/openid-connect/certs"
   AUTH_AUDIENCE: "essp-api"
   ```

---

## 4. Secret Management

### 4.1 Required Secrets List

#### Database Credentials
- `db-username`: PostgreSQL username
- `db-password`: PostgreSQL password
- Database DSN for each service (IMS, School, Devices, Parts)

#### Redis Credentials
- `redis-password`: Redis authentication password

#### NATS Credentials
- `nats-username`: NATS username (if auth enabled)
- `nats-password`: NATS password (if auth enabled)

#### MinIO Credentials
- `minio-access-key`: MinIO access key ID
- `minio-secret-key`: MinIO secret access key

#### JWT/Auth
- `jwt-secret`: JWT signing secret (32+ characters)

### 4.2 Creating Secrets in Kubernetes

#### Using kubectl

```bash
# Create namespace first
kubectl create namespace essp-prod

# Database secrets
kubectl create secret generic database-secrets \
  --namespace=essp-prod \
  --from-literal=db-username='essp_admin' \
  --from-literal=db-password='<strong-password>' \
  --from-literal=pg-dsn='postgres://essp_admin:<password>@postgresql:5432/ssp_ims_prod?sslmode=require'

# Redis secrets
kubectl create secret generic redis-secrets \
  --namespace=essp-prod \
  --from-literal=redis-password='<strong-redis-password>'

# MinIO secrets
kubectl create secret generic minio-secrets \
  --namespace=essp-prod \
  --from-literal=minio-access-key='essp_service_account' \
  --from-literal=minio-secret-key='<strong-minio-password>'

# NATS secrets (if authentication enabled)
kubectl create secret generic nats-secrets \
  --namespace=essp-prod \
  --from-literal=nats-username='essp_nats' \
  --from-literal=nats-password='<strong-nats-password>'

# JWT secret
kubectl create secret generic jwt-secrets \
  --namespace=essp-prod \
  --from-literal=jwt-secret="$(openssl rand -base64 32)"
```

#### Using Helm Values (Production)

For production, use `values-prod.yaml` with secret references:

```yaml
secrets:
  database:
    username: "essp_admin"
    password: "your-secure-password-here"  # Use external secret manager
  redis:
    password: "your-redis-password"
  minio:
    accessKey: "essp_service_account"
    secretKey: "your-minio-password"
  jwt:
    secret: "your-jwt-secret"
```

**DO NOT commit secrets to Git!** Use:
- Environment-specific values files (gitignored)
- External secret managers
- Encrypted secrets (sealed-secrets, SOPS)

#### Using Sealed Secrets (Recommended)

```bash
# Install sealed-secrets controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.5/controller.yaml

# Install kubeseal CLI
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.5/kubeseal-linux-amd64
chmod +x kubeseal-linux-amd64
sudo mv kubeseal-linux-amd64 /usr/local/bin/kubeseal

# Create sealed secret
kubectl create secret generic database-secrets \
  --namespace=essp-prod \
  --from-literal=db-password='<password>' \
  --dry-run=client -o yaml | \
  kubeseal -o yaml > sealed-database-secrets.yaml

# Apply sealed secret (safe to commit to Git)
kubectl apply -f sealed-database-secrets.yaml
```

#### Using HashiCorp Vault

```bash
# Store secrets in Vault
vault kv put secret/essp/prod/database \
  username=essp_admin \
  password=<strong-password> \
  dsn='postgres://...'

# Use Vault Agent Injector
# See: https://www.vaultproject.io/docs/platform/k8s/injector
```

### 4.3 Secret Rotation Procedures

#### Database Password Rotation

```bash
# Step 1: Update password in database
psql -h <db-host> -U essp_admin -d postgres -c "ALTER USER essp_admin PASSWORD '<new-password>';"

# Step 2: Update Kubernetes secret
kubectl create secret generic database-secrets \
  --namespace=essp-prod \
  --from-literal=db-password='<new-password>' \
  --from-literal=pg-dsn='postgres://essp_admin:<new-password>@...' \
  --dry-run=client -o yaml | kubectl apply -f -

# Step 3: Restart pods to pick up new secret
kubectl rollout restart deployment/ims-api -n essp-prod
kubectl rollout restart deployment/ssot-school -n essp-prod
kubectl rollout restart deployment/ssot-devices -n essp-prod
kubectl rollout restart deployment/ssot-parts -n essp-prod
kubectl rollout restart deployment/sync-worker -n essp-prod
```

#### Automated Rotation Schedule

| Secret Type | Rotation Frequency | Method |
|-------------|-------------------|--------|
| Database passwords | 90 days | Manual or automated via Vault |
| Redis passwords | 90 days | Manual rotation |
| MinIO access keys | 90 days | Create new key, update, delete old |
| JWT secrets | 180 days | Generate new, update, rolling restart |
| TLS certificates | 60 days before expiry | cert-manager auto-renewal |

---

## 5. Deployment Procedures

### 5.1 First-Time Deployment

#### Pre-Deployment Checklist

- [ ] Kubernetes cluster is running and accessible
- [ ] All infrastructure components deployed (PostgreSQL, Redis, NATS, MinIO, Keycloak)
- [ ] Secrets created in target namespace
- [ ] Container images built and pushed to registry
- [ ] DNS records configured for ingress
- [ ] TLS certificates ready or cert-manager configured
- [ ] Backup solution configured
- [ ] Monitoring stack deployed

#### Step-by-Step First Deployment

```bash
# 1. Clone repository
git clone <repository-url>
cd ESSP

# 2. Build and push container images
export REGISTRY=your-registry.example.com
export VERSION=v1.0.0

docker build -t $REGISTRY/essp/ims-api:$VERSION -f deployments/docker/ims-api.Dockerfile .
docker build -t $REGISTRY/essp/ssot-school:$VERSION -f deployments/docker/ssot-school.Dockerfile .
docker build -t $REGISTRY/essp/ssot-devices:$VERSION -f deployments/docker/ssot-devices.Dockerfile .
docker build -t $REGISTRY/essp/ssot-parts:$VERSION -f deployments/docker/ssot-parts.Dockerfile .
docker build -t $REGISTRY/essp/sync-worker:$VERSION -f deployments/docker/sync-worker.Dockerfile .

docker push $REGISTRY/essp/ims-api:$VERSION
docker push $REGISTRY/essp/ssot-school:$VERSION
docker push $REGISTRY/essp/ssot-devices:$VERSION
docker push $REGISTRY/essp/ssot-parts:$VERSION
docker push $REGISTRY/essp/sync-worker:$VERSION

# 3. Create namespace
kubectl create namespace essp-prod

# 4. Create image pull secret (if using private registry)
kubectl create secret docker-registry regcred \
  --docker-server=$REGISTRY \
  --docker-username=<username> \
  --docker-password=<password> \
  --namespace=essp-prod

# 5. Create secrets (see Section 4.2)
# [Execute secret creation commands from Section 4.2]

# 6. Run database migrations (see Section 6)
# [Execute migration commands from Section 6]

# 7. Deploy using Helm
helm install essp ./charts/essp \
  --namespace essp-prod \
  --create-namespace \
  --values charts/essp/values-prod.yaml \
  --set imsApi.image.tag=$VERSION \
  --set ssotSchool.image.tag=$VERSION \
  --set ssotDevices.image.tag=$VERSION \
  --set ssotParts.image.tag=$VERSION \
  --set syncWorker.image.tag=$VERSION \
  --set global.imageRegistry=$REGISTRY

# 8. Verify deployment
kubectl get pods -n essp-prod -w

# 9. Check rollout status
kubectl rollout status deployment/ims-api -n essp-prod
kubectl rollout status deployment/ssot-school -n essp-prod
kubectl rollout status deployment/ssot-devices -n essp-prod
kubectl rollout status deployment/ssot-parts -n essp-prod
kubectl rollout status deployment/sync-worker -n essp-prod

# 10. Test health endpoints
kubectl port-forward -n essp-prod svc/ims-api 8080:8080
curl http://localhost:8080/health
curl http://localhost:8080/readyz

# 11. Verify ingress
kubectl get ingress -n essp-prod
curl https://api.essp.example.com/health
```

### 5.2 Using Helm Chart

#### Install

```bash
helm install essp ./charts/essp \
  --namespace essp-prod \
  --create-namespace \
  --values charts/essp/values-prod.yaml
```

#### Upgrade

```bash
helm upgrade essp ./charts/essp \
  --namespace essp-prod \
  --values charts/essp/values-prod.yaml \
  --set imsApi.image.tag=v1.1.0
```

#### List Releases

```bash
helm list -n essp-prod
```

#### Get Values

```bash
# Get all values
helm get values essp -n essp-prod

# Get all values including defaults
helm get values essp -n essp-prod --all
```

#### History

```bash
helm history essp -n essp-prod
```

### 5.3 Using Raw Manifests

If you prefer not to use Helm:

```bash
# Apply in order
kubectl apply -f deployments/k8s/namespace.yaml
kubectl apply -f deployments/k8s/secrets/
kubectl apply -f deployments/k8s/configmaps/
kubectl apply -f deployments/k8s/services/
kubectl apply -f deployments/k8s/deployments/
kubectl apply -f deployments/k8s/ingress.yaml
kubectl apply -f deployments/k8s/pdbs/
kubectl apply -f deployments/k8s/hpa/
kubectl apply -f deployments/k8s/network-policies/
```

### 5.4 Environment-Specific Values

#### Development

```bash
helm install essp ./charts/essp \
  --namespace essp-dev \
  --create-namespace \
  --values charts/essp/values-dev.yaml
```

Key dev characteristics:
- Single replica per service
- Debug logging
- No authentication
- Lower resource limits
- Always pull images
- No HPA/PDB

#### Staging

```bash
helm install essp ./charts/essp \
  --namespace essp-staging \
  --create-namespace \
  --values charts/essp/values-staging.yaml
```

Key staging characteristics:
- 2 replicas for IMS API
- Info logging
- Authentication enabled
- Moderate resources
- HPA enabled
- Production-like configuration

#### Production

```bash
helm install essp ./charts/essp \
  --namespace essp-prod \
  --create-namespace \
  --values charts/essp/values-prod.yaml
```

Key production characteristics:
- 3+ replicas with HPA (up to 20)
- Warn logging
- Full authentication
- High resources
- PDB enabled
- Network policies enabled
- TLS required
- Pod anti-affinity

---

## 6. Database Migrations

### 6.1 Running Migrations

#### Migration Files Location

```
services/ims-api/migrations/
├── 001_init.sql
├── 002_bom.sql
├── 003_ssot_snapshots.sql
├── 004_ssot_checkpoints.sql
├── 005_service_delivery.sql
├── 006_work_orders_scheduling_deliverables.sql
├── 007_phase_checklists.sql
├── 008_audit_logs.sql
├── 900_ssot_snapshots.sql
├── 910_enrichment_denorm.sql
└── 920_work_order_parts_enrichment.sql
```

#### Method 1: Using Migration Job (Recommended)

Create a Kubernetes Job to run migrations:

```yaml
# migration-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: essp-migration-v1-0-0
  namespace: essp-prod
spec:
  ttlSecondsAfterFinished: 3600
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: migrate
        image: your-registry/essp/ims-api:v1.0.0
        command: ["/app/migrate"]
        args: ["up"]
        env:
        - name: PG_DSN
          valueFrom:
            secretKeyRef:
              name: database-secrets
              key: pg-dsn
```

Apply the job:

```bash
kubectl apply -f migration-job.yaml

# Watch job progress
kubectl logs -n essp-prod -f job/essp-migration-v1-0-0

# Verify completion
kubectl get job -n essp-prod essp-migration-v1-0-0
```

#### Method 2: Manual Migration

```bash
# Get database DSN from secret
export PG_DSN=$(kubectl get secret database-secrets -n essp-prod -o jsonpath='{.data.pg-dsn}' | base64 -d)

# Run migrations manually
kubectl run migrate-ims --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=your-registry/essp/ims-api:v1.0.0 \
  --env="PG_DSN=$PG_DSN" \
  -- /app/migrate up

# Or use psql directly
kubectl run psql-client --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=postgres:16 \
  --env="PGPASSWORD=<password>" \
  -- psql -h postgresql -U essp_admin -d ssp_ims_prod <<EOF
-- Paste migration SQL here
EOF
```

#### Method 3: Using golang-migrate

```bash
# Install golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Run migrations
migrate -path services/ims-api/migrations \
  -database "postgres://essp_admin:<password>@postgresql:5432/ssp_ims_prod?sslmode=require" \
  up
```

### 6.2 Rollback Procedures

#### Rollback Last Migration

```bash
# Using migration job
kubectl run migrate-rollback --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=your-registry/essp/ims-api:v1.0.0 \
  --env="PG_DSN=$PG_DSN" \
  -- /app/migrate down 1

# Using golang-migrate
migrate -path services/ims-api/migrations \
  -database "postgres://essp_admin:<password>@postgresql:5432/ssp_ims_prod" \
  down 1
```

#### Rollback to Specific Version

```bash
# Check current migration version
kubectl run migrate-version --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=your-registry/essp/ims-api:v1.0.0 \
  --env="PG_DSN=$PG_DSN" \
  -- /app/migrate version

# Rollback to version N
kubectl run migrate-rollback --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=your-registry/essp/ims-api:v1.0.0 \
  --env="PG_DSN=$PG_DSN" \
  -- /app/migrate goto N
```

### 6.3 Schema Validation

#### Verify Migration Status

```bash
# Check schema_migrations table
kubectl run psql-client --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=postgres:16 \
  --env="PGPASSWORD=<password>" \
  -- psql -h postgresql -U essp_admin -d ssp_ims_prod -c "SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 10;"
```

#### Validate Table Structure

```bash
# List all tables
kubectl run psql-client --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=postgres:16 \
  --env="PGPASSWORD=<password>" \
  -- psql -h postgresql -U essp_admin -d ssp_ims_prod -c "\dt"

# Describe specific table
kubectl run psql-client --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=postgres:16 \
  --env="PGPASSWORD=<password>" \
  -- psql -h postgresql -U essp_admin -d ssp_ims_prod -c "\d incidents"
```

#### Pre-Deployment Migration Checklist

- [ ] Backup database before running migrations
- [ ] Test migrations in staging environment first
- [ ] Review migration SQL for destructive operations
- [ ] Ensure migrations are idempotent
- [ ] Verify rollback scripts exist
- [ ] Check for blocking locks or long-running queries
- [ ] Plan for maintenance window if needed
- [ ] Monitor migration progress
- [ ] Validate data integrity after migration

---

## 7. Configuration

### 7.1 ConfigMaps

#### IMS API Configuration

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ims-api-config
  namespace: essp-prod
data:
  APP_ENV: "production"
  HTTP_ADDR: ":8080"
  LOG_LEVEL: "warn"

  # Database
  # Note: DSN is in secrets, not configmap

  # Redis
  REDIS_ADDR: "redis-master.essp-prod.svc.cluster.local:6379"
  REDIS_DB: "0"

  # NATS
  NATS_URL: "nats://nats.essp-prod.svc.cluster.local:4222"

  # MinIO
  MINIO_ENDPOINT: "minio.essp-prod.svc.cluster.local:9000"
  MINIO_USE_SSL: "true"
  MINIO_REGION: "us-east-1"
  MINIO_PRESIGN_EXPIRY_SECONDS: "3600"
  ATTACHMENTS_BUCKET: "essp-attachments-prod"
  ATTACHMENTS_PUBLIC_BASE_URL: "https://attachments.essp.example.com"

  # SSOT Services
  SCHOOL_SSOT_BASE_URL: "http://ssot-school.essp-prod.svc.cluster.local:8081"
  DEVICE_SSOT_BASE_URL: "http://ssot-devices.essp-prod.svc.cluster.local:8082"
  PARTS_SSOT_BASE_URL: "http://ssot-parts.essp-prod.svc.cluster.local:8083"
  SSOT_SYNC_PAGE_SIZE: "100"

  # Authentication
  AUTH_ENABLED: "true"
  AUTH_AUDIENCE: "essp-api"
  AUTH_ISSUER: "https://auth.essp.example.com/realms/essp"
  AUTH_JWKS_URL: "https://auth.essp.example.com/realms/essp/protocol/openid-connect/certs"

  # Tenancy
  TENANT_HEADER: "X-Tenant-ID"
  SCHOOL_HEADER: "X-School-ID"

  # CORS
  CORS_ALLOWED_ORIGINS: "https://app.essp.example.com,https://admin.essp.example.com"

  # Work Orders
  AUTO_ROUTE_WORK_ORDERS: "true"
  DEFAULT_REPAIR_LOCATION: "central-depot"

  # Rate Limiting
  RATE_LIMIT_ENABLED: "true"
  RATE_LIMIT_READ_RPM: "1000"
  RATE_LIMIT_WRITE_RPM: "100"
  RATE_LIMIT_BURST: "50"
```

#### Updating ConfigMaps

```bash
# Edit configmap
kubectl edit configmap ims-api-config -n essp-prod

# Or apply updated file
kubectl apply -f deployments/k8s/configmaps/ims-api-config.yaml

# Restart pods to pick up changes
kubectl rollout restart deployment/ims-api -n essp-prod
```

### 7.2 Environment Variables

#### Precedence Order

1. Container-level env vars (highest priority)
2. ConfigMap values
3. Secret values
4. Default application values (lowest priority)

#### Common Environment Variables

| Variable | Source | Example Value | Description |
|----------|--------|---------------|-------------|
| `APP_ENV` | ConfigMap | production | Environment name |
| `HTTP_ADDR` | ConfigMap | :8080 | HTTP listen address |
| `LOG_LEVEL` | ConfigMap | warn | Logging level |
| `PG_DSN` | Secret | postgres://... | Database connection string |
| `REDIS_ADDR` | ConfigMap | redis:6379 | Redis address |
| `REDIS_PASSWORD` | Secret | *** | Redis password |
| `NATS_URL` | ConfigMap | nats://nats:4222 | NATS connection URL |
| `MINIO_ENDPOINT` | ConfigMap | minio:9000 | MinIO endpoint |
| `MINIO_ACCESS_KEY` | Secret | *** | MinIO access key |
| `MINIO_SECRET_KEY` | Secret | *** | MinIO secret key |
| `AUTH_ENABLED` | ConfigMap | true | Enable authentication |
| `AUTH_JWKS_URL` | ConfigMap | https://... | JWKS endpoint |

### 7.3 Feature Flags

Feature flags can be controlled via ConfigMap:

```yaml
# Feature flags in ConfigMap
data:
  # Authentication
  AUTH_ENABLED: "true"

  # CORS
  CORS_ENABLED: "true"

  # Rate Limiting
  RATE_LIMIT_ENABLED: "true"

  # Work Order Features
  AUTO_ROUTE_WORK_ORDERS: "true"
  ENABLE_WORK_ORDER_SCHEDULING: "true"

  # SSOT Sync
  ENABLE_SSOT_SYNC: "true"
  SSOT_SYNC_INTERVAL: "5m"
```

#### Toggling Features

```bash
# Disable authentication (dev/testing only)
kubectl patch configmap ims-api-config -n essp-dev \
  --type merge \
  -p '{"data":{"AUTH_ENABLED":"false"}}'

# Restart to apply
kubectl rollout restart deployment/ims-api -n essp-dev

# Enable auto-routing
kubectl patch configmap ims-api-config -n essp-prod \
  --type merge \
  -p '{"data":{"AUTO_ROUTE_WORK_ORDERS":"true"}}'

kubectl rollout restart deployment/ims-api -n essp-prod
```

---

## 8. Health Checks

### 8.1 Readiness Probes

Readiness probes determine if a pod should receive traffic.

#### Configuration

```yaml
readinessProbe:
  httpGet:
    path: /readyz
    port: http
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
  successThreshold: 1
```

#### /readyz Endpoint

Checks:
- Database connectivity
- Redis connectivity
- NATS connectivity
- MinIO connectivity
- SSOT service connectivity (for IMS API)

#### Testing Readiness

```bash
# Port-forward to service
kubectl port-forward -n essp-prod svc/ims-api 8080:8080

# Test readiness endpoint
curl http://localhost:8080/readyz

# Expected response (200 OK)
{
  "status": "ready",
  "checks": {
    "database": "ok",
    "redis": "ok",
    "nats": "ok",
    "minio": "ok"
  }
}
```

### 8.2 Liveness Probes

Liveness probes determine if a pod should be restarted.

#### Configuration

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

#### /health Endpoint

Simple health check that returns 200 if the service is alive.

```bash
curl http://localhost:8080/health

# Expected response (200 OK)
{
  "status": "ok"
}
```

### 8.3 Startup Probes

Startup probes provide additional time for slow-starting applications.

#### Configuration (Optional)

```yaml
startupProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 0
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 30  # Allow 150 seconds for startup
```

Useful for:
- Applications with slow initialization
- Services that need to warm up caches
- Migration-heavy deployments

---

## 9. Scaling

### 9.1 Horizontal Pod Autoscaler (HPA)

#### IMS API HPA Configuration

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ims-api-hpa
  namespace: essp-prod
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ims-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 15
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
      - type: Pods
        value: 2
        periodSeconds: 15
      selectPolicy: Max
```

#### Monitoring HPA

```bash
# Get HPA status
kubectl get hpa -n essp-prod

# Describe HPA
kubectl describe hpa ims-api-hpa -n essp-prod

# Watch HPA
kubectl get hpa -n essp-prod -w

# Check metrics
kubectl top pods -n essp-prod -l app=ims-api
```

#### Custom Metrics (Advanced)

Using custom metrics like request rate:

```yaml
metrics:
- type: Pods
  pods:
    metric:
      name: http_requests_per_second
    target:
      type: AverageValue
      averageValue: "1000"
```

### 9.2 Manual Scaling

#### Scale Up

```bash
# Scale IMS API to 5 replicas
kubectl scale deployment ims-api -n essp-prod --replicas=5

# Scale all SSOT services
kubectl scale deployment ssot-school ssot-devices ssot-parts -n essp-prod --replicas=2

# Verify scaling
kubectl get deployment -n essp-prod
```

#### Scale Down

```bash
# Scale down during maintenance
kubectl scale deployment ims-api -n essp-prod --replicas=1

# Scale to zero (stop service)
kubectl scale deployment sync-worker -n essp-prod --replicas=0
```

### 9.3 Resource Tuning

#### Identify Resource Usage

```bash
# Current usage
kubectl top pods -n essp-prod

# Detailed metrics
kubectl describe nodes | grep -A 5 "Allocated resources"
```

#### Right-Sizing Resources

```yaml
# Conservative (default)
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi

# Medium load
resources:
  requests:
    cpu: 200m
    memory: 256Mi
  limits:
    cpu: 1000m
    memory: 1Gi

# High load
resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 2000m
    memory: 2Gi
```

#### Vertical Pod Autoscaler (VPA)

```bash
# Install VPA
kubectl apply -f https://github.com/kubernetes/autoscaler/releases/download/vertical-pod-autoscaler-0.14.0/vpa-v0.14.0.yaml

# Create VPA
cat <<EOF | kubectl apply -f -
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: ims-api-vpa
  namespace: essp-prod
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ims-api
  updatePolicy:
    updateMode: "Auto"
EOF

# Check VPA recommendations
kubectl describe vpa ims-api-vpa -n essp-prod
```

---

## 10. Monitoring & Alerting

### 10.1 Prometheus Metrics

#### Metrics Endpoints

All ESSP services expose Prometheus metrics at `/metrics`:

```bash
# Port-forward to service
kubectl port-forward -n essp-prod svc/ims-api 8080:8080

# Scrape metrics
curl http://localhost:8080/metrics
```

#### Key Metrics

**HTTP Metrics**:
- `http_requests_total` - Total HTTP requests (labels: method, path, status)
- `http_request_duration_seconds` - Request duration histogram
- `http_requests_in_flight` - Current in-flight requests

**Database Metrics**:
- `db_connections_active` - Active database connections
- `db_query_duration_seconds` - Query duration histogram

**Business Metrics**:
- `incidents_created_total` - Total incidents created
- `work_orders_created_total` - Total work orders created

#### ServiceMonitor Configuration

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: essp-services
  namespace: essp-prod
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app.kubernetes.io/part-of: essp-platform
  endpoints:
  - port: http
    path: /metrics
    interval: 15s
    scrapeTimeout: 10s
```

Deploy ServiceMonitor:

```bash
kubectl apply -f deployments/monitoring/servicemonitor.yaml
```

### 10.2 Grafana Dashboards

#### Access Grafana

```bash
# Port-forward to Grafana
kubectl port-forward -n monitoring svc/grafana 3000:3000

# Open http://localhost:3000
# Default credentials: admin/admin
```

#### Available Dashboards

1. **ESSP Platform Overview** (UID: `essp-overview`)
   - Platform availability
   - Request rate by service
   - Error rate
   - Latency percentiles
   - Database connections

2. **IMS API Metrics** (UID: `ims-api`)
   - Request rate by method
   - Response status distribution
   - Latency percentiles
   - Top slowest endpoints
   - Business metrics

#### Importing Dashboards

```bash
# Dashboards are in deployments/monitoring/grafana/dashboards/

# Import via UI
# 1. Navigate to Grafana
# 2. Click + → Import
# 3. Upload JSON file

# Or provision automatically via configmap
kubectl create configmap grafana-dashboards \
  --from-file=deployments/monitoring/grafana/dashboards/ \
  --namespace=monitoring
```

### 10.3 AlertManager Rules

#### Critical Alerts

```yaml
# From deployments/monitoring/prometheus/alerts.yml

groups:
  - name: essp_api_alerts
    rules:
      - alert: HighErrorRate
        expr: |
          (sum(rate(http_requests_total{status=~"5.."}[5m])) by (service) /
           sum(rate(http_requests_total[5m])) by (service)) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate in {{ $labels.service }}"

      - alert: ServiceDown
        expr: up{job=~"ims-api|ssot-.*|sync-worker"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.job }} is down"
```

#### Configure Notifications

Edit `deployments/monitoring/alertmanager/config.yml`:

```yaml
route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default-receiver'
  routes:
  - match:
      severity: critical
    receiver: critical-receiver

receivers:
  - name: 'default-receiver'
    email_configs:
      - to: 'team@essp.example.com'

  - name: 'critical-receiver'
    pagerduty_configs:
      - service_key: '<pagerduty-key>'
    slack_configs:
      - api_url: '<slack-webhook-url>'
        channel: '#essp-critical-alerts'
```

#### Deploy AlertManager

```bash
# Apply alertmanager configuration
kubectl create configmap alertmanager-config \
  --from-file=deployments/monitoring/alertmanager/config.yml \
  --namespace=monitoring \
  --dry-run=client -o yaml | kubectl apply -f -

# Reload AlertManager
kubectl rollout restart statefulset/alertmanager -n monitoring
```

### 10.4 Log Aggregation

#### Using Loki + Promtail

```bash
# Add Grafana Loki repo
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Install Loki
helm install loki grafana/loki-stack \
  --namespace monitoring \
  --set promtail.enabled=true \
  --set loki.persistence.enabled=true \
  --set loki.persistence.size=10Gi

# Verify installation
kubectl get pods -n monitoring -l app=loki
kubectl get pods -n monitoring -l app=promtail
```

#### Query Logs in Grafana

```bash
# Add Loki as datasource in Grafana
# URL: http://loki:3100

# Example LogQL queries:
# All logs from ims-api
{namespace="essp-prod", app="ims-api"}

# Error logs
{namespace="essp-prod"} |= "error"

# Logs from specific pod
{namespace="essp-prod", pod="ims-api-xxxxx-xxxxx"}
```

#### ELK Stack Alternative

For ELK (Elasticsearch, Logstash, Kibana):

```bash
# Install ECK Operator
kubectl apply -f https://download.elastic.co/downloads/eck/2.10.0/crds.yaml
kubectl apply -f https://download.elastic.co/downloads/eck/2.10.0/operator.yaml

# Deploy Elasticsearch
kubectl apply -f - <<EOF
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: essp-logs
  namespace: monitoring
spec:
  version: 8.11.0
  nodeSets:
  - name: default
    count: 3
    config:
      node.store.allow_mmap: false
EOF
```

---

## 11. Troubleshooting

### 11.1 Common Issues

#### Pods Not Starting

**Symptoms**:
- Pods stuck in `Pending`, `ContainerCreating`, or `CrashLoopBackOff`

**Diagnosis**:
```bash
# Check pod status
kubectl get pods -n essp-prod

# Describe pod for events
kubectl describe pod <pod-name> -n essp-prod

# Check logs
kubectl logs <pod-name> -n essp-prod
kubectl logs <pod-name> -n essp-prod --previous  # Previous container logs
```

**Common Causes**:

1. **Image Pull Errors**
   ```bash
   # Check image pull secret
   kubectl get secret regcred -n essp-prod

   # Verify image exists
   docker pull <image>
   ```

2. **Insufficient Resources**
   ```bash
   # Check node capacity
   kubectl describe nodes | grep -A 5 "Allocated resources"

   # Reduce resource requests or add nodes
   ```

3. **Missing Secrets/ConfigMaps**
   ```bash
   # List secrets
   kubectl get secrets -n essp-prod

   # List configmaps
   kubectl get configmaps -n essp-prod
   ```

4. **Failed Health Checks**
   ```bash
   # Check readiness probe
   kubectl describe pod <pod-name> -n essp-prod | grep -A 10 "Readiness"

   # Test endpoint manually
   kubectl port-forward <pod-name> 8080:8080 -n essp-prod
   curl http://localhost:8080/readyz
   ```

#### Database Connection Failures

**Symptoms**:
- Logs show "connection refused", "timeout", or "authentication failed"

**Diagnosis**:
```bash
# Test database connection from pod
kubectl run psql-test --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=postgres:16 \
  --env="PGPASSWORD=<password>" \
  -- psql -h <db-host> -U essp_admin -d ssp_ims_prod -c "SELECT 1;"

# Check database secret
kubectl get secret database-secrets -n essp-prod -o yaml
```

**Solutions**:
1. Verify database is running and accessible
2. Check credentials in secrets
3. Verify network policies allow egress to database
4. Check database firewall rules

#### Service Not Accessible via Ingress

**Symptoms**:
- 404, 502, or 503 errors when accessing via ingress

**Diagnosis**:
```bash
# Check ingress
kubectl get ingress -n essp-prod
kubectl describe ingress essp-ingress -n essp-prod

# Check service endpoints
kubectl get endpoints -n essp-prod

# Check ingress controller logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
```

**Solutions**:
1. Verify pods are ready
   ```bash
   kubectl get pods -n essp-prod -l app=ims-api
   ```

2. Check service selector matches pod labels
   ```bash
   kubectl get svc ims-api -n essp-prod -o yaml
   kubectl get pods -n essp-prod -l app=ims-api --show-labels
   ```

3. Verify ingress annotations
4. Check DNS records point to ingress controller

### 11.2 Debug Commands

#### Pod Debugging

```bash
# Get all pods with status
kubectl get pods -n essp-prod -o wide

# Filter by status
kubectl get pods -n essp-prod --field-selector=status.phase!=Running

# Describe pod
kubectl describe pod <pod-name> -n essp-prod

# Get pod logs
kubectl logs <pod-name> -n essp-prod
kubectl logs <pod-name> -n essp-prod -f  # Follow
kubectl logs <pod-name> -n essp-prod --tail=100  # Last 100 lines
kubectl logs <pod-name> -n essp-prod --since=1h  # Last hour

# Exec into pod
kubectl exec -it <pod-name> -n essp-prod -- /bin/sh

# Port-forward
kubectl port-forward <pod-name> 8080:8080 -n essp-prod

# Copy files from pod
kubectl cp essp-prod/<pod-name>:/path/to/file ./local-file

# Get pod YAML
kubectl get pod <pod-name> -n essp-prod -o yaml
```

#### Deployment Debugging

```bash
# Get deployments
kubectl get deployments -n essp-prod

# Describe deployment
kubectl describe deployment ims-api -n essp-prod

# Check rollout status
kubectl rollout status deployment/ims-api -n essp-prod

# Check rollout history
kubectl rollout history deployment/ims-api -n essp-prod

# Check replica sets
kubectl get rs -n essp-prod
```

#### Service Debugging

```bash
# Get services
kubectl get svc -n essp-prod

# Describe service
kubectl describe svc ims-api -n essp-prod

# Check endpoints
kubectl get endpoints ims-api -n essp-prod

# Test service from within cluster
kubectl run curl-test --rm -it --restart=Never \
  --image=curlimages/curl:latest \
  -- curl http://ims-api.essp-prod.svc.cluster.local:8080/health
```

#### Network Debugging

```bash
# Check network policies
kubectl get networkpolicies -n essp-prod

# Describe network policy
kubectl describe networkpolicy <policy-name> -n essp-prod

# Test DNS resolution
kubectl run dnsutils --rm -it --restart=Never \
  --image=gcr.io/kubernetes-e2e-test-images/dnsutils:1.3 \
  -- nslookup ims-api.essp-prod.svc.cluster.local

# Test connectivity
kubectl run netshoot --rm -it --restart=Never \
  --image=nicolaka/netshoot \
  -- curl http://ims-api.essp-prod.svc.cluster.local:8080/health
```

### 11.3 Log Analysis

#### View Application Logs

```bash
# All logs from IMS API
kubectl logs -n essp-prod -l app=ims-api --tail=100

# All logs from all ESSP services
kubectl logs -n essp-prod -l app.kubernetes.io/part-of=essp-platform --tail=50

# Follow logs in real-time
kubectl logs -n essp-prod -l app=ims-api -f

# Search logs for errors
kubectl logs -n essp-prod -l app=ims-api --tail=1000 | grep -i error

# Get logs from all pods in deployment
for pod in $(kubectl get pods -n essp-prod -l app=ims-api -o name); do
  echo "=== $pod ==="
  kubectl logs -n essp-prod $pod --tail=20
done
```

#### Structured Log Parsing

```bash
# Parse JSON logs with jq
kubectl logs -n essp-prod <pod-name> | jq 'select(.level == "error")'

# Count errors by type
kubectl logs -n essp-prod <pod-name> | jq -r '.error_type' | sort | uniq -c

# Filter by timestamp
kubectl logs -n essp-prod <pod-name> | jq 'select(.timestamp > "2025-12-12T10:00:00Z")'
```

#### Event Logs

```bash
# Get events in namespace
kubectl get events -n essp-prod --sort-by='.lastTimestamp'

# Watch events
kubectl get events -n essp-prod -w

# Filter events by type
kubectl get events -n essp-prod --field-selector type=Warning

# Events for specific object
kubectl get events -n essp-prod --field-selector involvedObject.name=ims-api
```

---

## 12. Rollback Procedures

### 12.1 Helm Rollback

#### List Revisions

```bash
# Show deployment history
helm history essp -n essp-prod

# Output:
# REVISION  UPDATED                   STATUS      CHART         APP VERSION  DESCRIPTION
# 1         Mon Dec 11 10:00:00 2025  superseded  essp-0.1.0    1.0.0        Install complete
# 2         Mon Dec 12 14:00:00 2025  deployed    essp-0.1.0    1.1.0        Upgrade complete
```

#### Rollback to Previous Version

```bash
# Rollback to previous revision
helm rollback essp -n essp-prod

# Rollback to specific revision
helm rollback essp 1 -n essp-prod

# Rollback with wait
helm rollback essp -n essp-prod --wait --timeout=10m

# Dry-run rollback
helm rollback essp -n essp-prod --dry-run
```

#### Verify Rollback

```bash
# Check rollout status
kubectl rollout status deployment/ims-api -n essp-prod

# Verify version
kubectl get deployment ims-api -n essp-prod -o jsonpath='{.spec.template.spec.containers[0].image}'

# Test application
curl https://api.essp.example.com/health
```

### 12.2 Manual Rollback

#### Rollback Deployment

```bash
# View rollout history
kubectl rollout history deployment/ims-api -n essp-prod

# Rollback to previous version
kubectl rollout undo deployment/ims-api -n essp-prod

# Rollback to specific revision
kubectl rollout undo deployment/ims-api -n essp-prod --to-revision=3

# Check rollback status
kubectl rollout status deployment/ims-api -n essp-prod
```

#### Rollback Multiple Services

```bash
# Rollback all ESSP services
for deployment in ims-api ssot-school ssot-devices ssot-parts sync-worker; do
  echo "Rolling back $deployment..."
  kubectl rollout undo deployment/$deployment -n essp-prod
  kubectl rollout status deployment/$deployment -n essp-prod
done
```

### 12.3 Database Rollback

#### Rollback Migration

```bash
# Check current migration version
kubectl run migrate-check --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=your-registry/essp/ims-api:v1.0.0 \
  --env="PG_DSN=$PG_DSN" \
  -- /app/migrate version

# Rollback to previous migration
kubectl run migrate-rollback --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=your-registry/essp/ims-api:v1.0.0 \
  --env="PG_DSN=$PG_DSN" \
  -- /app/migrate down 1

# Rollback to specific version
kubectl run migrate-rollback --rm -it --restart=Never \
  --namespace=essp-prod \
  --image=your-registry/essp/ims-api:v1.0.0 \
  --env="PG_DSN=$PG_DSN" \
  -- /app/migrate goto 5
```

#### Restore from Backup

```bash
# Stop application first
kubectl scale deployment ims-api ssot-school ssot-devices ssot-parts sync-worker \
  -n essp-prod --replicas=0

# Restore database from backup
pg_restore -h <db-host> -U essp_admin -d ssp_ims_prod \
  --clean --if-exists backup_file.dump

# Or for SQL dump
psql -h <db-host> -U essp_admin -d ssp_ims_prod < backup_file.sql

# Restart application
kubectl scale deployment ims-api -n essp-prod --replicas=3
kubectl scale deployment ssot-school ssot-devices ssot-parts -n essp-prod --replicas=2
kubectl scale deployment sync-worker -n essp-prod --replicas=1
```

#### Rollback Decision Matrix

| Scenario | Action | Estimated Time |
|----------|--------|----------------|
| Minor bug in code | Helm/kubectl rollback | 2-5 minutes |
| Configuration error | Update ConfigMap/Secret + restart | 5-10 minutes |
| Database schema issue | Migration rollback | 10-30 minutes |
| Data corruption | Restore from backup | 30-120 minutes |
| Complete system failure | Full stack rollback + DB restore | 60-180 minutes |

---

## 13. Disaster Recovery

### 13.1 Backup Procedures

#### Database Backups

**Automated Backups (Recommended)**:

```bash
# Create CronJob for daily backups
cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: CronJob
metadata:
  name: postgres-backup
  namespace: essp-prod
spec:
  schedule: "0 2 * * *"  # 2 AM daily
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:16
            env:
            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: database-secrets
                  key: db-password
            command:
            - /bin/sh
            - -c
            - |
              TIMESTAMP=\$(date +%Y%m%d_%H%M%S)
              pg_dump -h postgresql -U essp_admin -F c ssp_ims_prod > /backup/ssp_ims_prod_\${TIMESTAMP}.dump
              # Upload to S3 or other storage
              aws s3 cp /backup/ssp_ims_prod_\${TIMESTAMP}.dump s3://essp-backups/database/
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
          volumes:
          - name: backup-storage
            persistentVolumeClaim:
              claimName: backup-pvc
          restartPolicy: OnFailure
EOF
```

**Manual Backup**:

```bash
# Backup all databases
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# IMS database
pg_dump -h <db-host> -U essp_admin -F c ssp_ims_prod > ssp_ims_prod_${TIMESTAMP}.dump

# SSOT databases
pg_dump -h <db-host> -U essp_admin -F c ssp_school_prod > ssp_school_prod_${TIMESTAMP}.dump
pg_dump -h <db-host> -U essp_admin -F c ssp_devices_prod > ssp_devices_prod_${TIMESTAMP}.dump
pg_dump -h <db-host> -U essp_admin -F c ssp_parts_prod > ssp_parts_prod_${TIMESTAMP}.dump

# Upload to S3
aws s3 cp ssp_ims_prod_${TIMESTAMP}.dump s3://essp-backups/database/
```

#### MinIO/S3 Backups

```bash
# Sync MinIO bucket to S3 for backup
mc mirror essp/essp-attachments-prod s3-backup/essp-attachments-prod

# Or use rclone
rclone sync minio:essp-attachments-prod s3:essp-backups/attachments/
```

#### Kubernetes Configuration Backup

```bash
# Backup all Kubernetes resources
kubectl get all,cm,secret,ing,pdb,hpa,networkpolicy \
  -n essp-prod \
  -o yaml > essp-prod-backup-$(date +%Y%m%d).yaml

# Backup Helm values
helm get values essp -n essp-prod > essp-helm-values-$(date +%Y%m%d).yaml
```

#### Velero for Cluster Backups

```bash
# Install Velero
wget https://github.com/vmware-tanzu/velero/releases/download/v1.12.1/velero-v1.12.1-linux-amd64.tar.gz
tar -xvf velero-v1.12.1-linux-amd64.tar.gz
sudo mv velero-v1.12.1-linux-amd64/velero /usr/local/bin/

# Install Velero in cluster (AWS example)
velero install \
  --provider aws \
  --plugins velero/velero-plugin-for-aws:v1.8.0 \
  --bucket essp-velero-backups \
  --secret-file ./credentials-velero \
  --backup-location-config region=us-east-1 \
  --snapshot-location-config region=us-east-1

# Create backup schedule
velero schedule create essp-daily \
  --schedule="0 2 * * *" \
  --include-namespaces essp-prod \
  --ttl 720h

# Create manual backup
velero backup create essp-backup-$(date +%Y%m%d) \
  --include-namespaces essp-prod
```

### 13.2 Restore Procedures

#### Database Restore

```bash
# Stop application pods
kubectl scale deployment ims-api ssot-school ssot-devices ssot-parts sync-worker \
  -n essp-prod --replicas=0

# Restore database
pg_restore -h <db-host> -U essp_admin \
  --clean --if-exists \
  -d ssp_ims_prod \
  ssp_ims_prod_20251212_020000.dump

# Verify data
psql -h <db-host> -U essp_admin -d ssp_ims_prod -c "SELECT COUNT(*) FROM incidents;"

# Restart application
kubectl scale deployment ims-api -n essp-prod --replicas=3
kubectl scale deployment ssot-school ssot-devices ssot-parts -n essp-prod --replicas=2
kubectl scale deployment sync-worker -n essp-prod --replicas=1
```

#### Point-in-Time Recovery (PITR)

```bash
# For PostgreSQL with WAL archiving
pg_basebackup -h <db-host> -U essp_admin -D /backup/base -Fp -Xs -P

# Restore to specific point in time
# 1. Stop database
# 2. Restore base backup
# 3. Configure recovery.conf with target time
# 4. Start database
```

#### Kubernetes Resource Restore

```bash
# Restore from backup file
kubectl apply -f essp-prod-backup-20251212.yaml

# Or use Velero
velero restore create --from-backup essp-backup-20251212

# Check restore status
velero restore describe essp-backup-20251212-restore
```

### 13.3 RTO/RPO Targets

#### Recovery Objectives

| Environment | RTO (Recovery Time Objective) | RPO (Recovery Point Objective) | Backup Frequency |
|-------------|-------------------------------|--------------------------------|------------------|
| **Production** | 4 hours | 1 hour | Continuous (WAL) + Daily full |
| **Staging** | 8 hours | 24 hours | Daily |
| **Development** | 24 hours | 1 week | Weekly |

#### RTO Components

| Task | Estimated Time | Notes |
|------|---------------|-------|
| Detect incident | 5-15 minutes | Automated monitoring |
| Assess damage | 15-30 minutes | Team evaluation |
| Retrieve backup | 10-30 minutes | From S3/backup storage |
| Restore database | 30-120 minutes | Depends on DB size |
| Restore Kubernetes resources | 10-20 minutes | kubectl apply |
| Verify and test | 30-60 minutes | End-to-end testing |
| Total RTO | ~2-4 hours | Production target |

#### DR Checklist

- [ ] Database backups verified and tested
- [ ] MinIO/S3 data replicated
- [ ] Kubernetes manifests backed up
- [ ] Secrets stored in secure backup location
- [ ] DNS records documented
- [ ] TLS certificates backed up
- [ ] Disaster recovery plan documented
- [ ] Recovery procedures tested quarterly
- [ ] Team trained on DR procedures
- [ ] Contact list updated
- [ ] Alternative infrastructure available
- [ ] Communication plan in place

---

## 14. Maintenance

### 14.1 Scheduled Maintenance Windows

#### Recommended Windows

| Environment | Window | Duration | Frequency |
|-------------|--------|----------|-----------|
| **Production** | Sunday 2:00-6:00 AM UTC | 4 hours | Monthly |
| **Staging** | Saturday 2:00-4:00 AM UTC | 2 hours | Bi-weekly |
| **Development** | Any time | N/A | As needed |

#### Maintenance Procedure

```bash
# 1. Notify users (production only)
# Send notification via email/Slack 48 hours in advance

# 2. Create backup
# See Section 13.1

# 3. Put application in maintenance mode (optional)
kubectl patch ingress essp-ingress -n essp-prod \
  --type=json \
  -p='[{"op": "add", "path": "/metadata/annotations/nginx.ingress.kubernetes.io~1default-backend", "value": "maintenance-backend"}]'

# 4. Scale down to minimum replicas
kubectl scale deployment ims-api -n essp-prod --replicas=1
kubectl scale deployment ssot-school ssot-devices ssot-parts -n essp-prod --replicas=1

# 5. Perform maintenance
# - Apply updates
# - Run migrations
# - Update configurations

# 6. Scale back up
kubectl scale deployment ims-api -n essp-prod --replicas=3
kubectl scale deployment ssot-school ssot-devices ssot-parts -n essp-prod --replicas=2

# 7. Verify health
kubectl get pods -n essp-prod
curl https://api.essp.example.com/health

# 8. Remove maintenance mode
kubectl patch ingress essp-ingress -n essp-prod --type=json \
  -p='[{"op": "remove", "path": "/metadata/annotations/nginx.ingress.kubernetes.io~1default-backend"}]'

# 9. Monitor for issues
kubectl logs -n essp-prod -l app=ims-api -f
```

### 14.2 Certificate Renewal

#### Using cert-manager (Automated)

```bash
# Check certificate status
kubectl get certificates -n essp-prod

# Describe certificate
kubectl describe certificate essp-prod-tls -n essp-prod

# Force renewal
cmctl renew essp-prod-tls -n essp-prod

# Check certificate expiry
kubectl get secret essp-prod-tls -n essp-prod -o jsonpath='{.data.tls\.crt}' | \
  base64 -d | openssl x509 -noout -enddate
```

#### Manual Certificate Renewal

```bash
# Generate new certificate using Let's Encrypt
certbot certonly --manual \
  --preferred-challenges dns \
  -d api.essp.example.com

# Update Kubernetes secret
kubectl create secret tls essp-prod-tls \
  --cert=/etc/letsencrypt/live/api.essp.example.com/fullchain.pem \
  --key=/etc/letsencrypt/live/api.essp.example.com/privkey.pem \
  --namespace=essp-prod \
  --dry-run=client -o yaml | kubectl apply -f -

# Restart ingress controller to pick up new cert
kubectl rollout restart deployment ingress-nginx-controller -n ingress-nginx
```

#### Certificate Expiry Monitoring

```bash
# Create alert for certificate expiry
# Add to prometheus/alerts.yml
- alert: CertificateExpiringSoon
  expr: |
    (x509_cert_not_after - time()) / 86400 < 30
  for: 1h
  labels:
    severity: warning
  annotations:
    summary: "Certificate {{ $labels.subject_CN }} expiring soon"
    description: "Certificate expires in {{ $value }} days"
```

### 14.3 Dependency Updates

#### Kubernetes Cluster Updates

```bash
# Check current version
kubectl version

# Update control plane (managed cluster)
# AWS EKS
aws eks update-cluster-version --name essp-cluster --kubernetes-version 1.28

# GKE
gcloud container clusters upgrade essp-cluster --master --cluster-version 1.28

# Update node pools
# EKS
aws eks update-nodegroup-version --cluster-name essp-cluster --nodegroup-name essp-nodes

# GKE
gcloud container clusters upgrade essp-cluster --node-pool=default-pool
```

#### Application Dependencies

```bash
# Update Go dependencies
cd services/ims-api
go get -u ./...
go mod tidy

# Rebuild and test
go test ./...

# Build new image
docker build -t your-registry/essp/ims-api:v1.1.0 .
docker push your-registry/essp/ims-api:v1.1.0
```

#### Infrastructure Updates

```bash
# Update PostgreSQL version
# 1. Create read replica with new version
# 2. Test application compatibility
# 3. Promote replica to primary
# 4. Update connection string

# Update Redis
helm upgrade redis bitnami/redis \
  --namespace essp-prod \
  --reuse-values \
  --set image.tag=7.2

# Update NATS
helm upgrade nats nats/nats \
  --namespace essp-prod \
  --reuse-values \
  --version 1.0.0

# Update MinIO
helm upgrade minio minio/minio \
  --namespace essp-prod \
  --reuse-values \
  --set image.tag=RELEASE.2023-12-09T18-17-51Z
```

#### Security Patches

```bash
# Scan images for vulnerabilities
trivy image your-registry/essp/ims-api:v1.0.0

# Update base images
# Edit Dockerfile
FROM golang:1.21 AS builder  # Update to 1.21.5

# Rebuild and deploy
docker build -t your-registry/essp/ims-api:v1.0.1 .
docker push your-registry/essp/ims-api:v1.0.1

helm upgrade essp ./charts/essp \
  --namespace essp-prod \
  --reuse-values \
  --set imsApi.image.tag=v1.0.1
```

#### Update Schedule

| Component | Check Frequency | Update Frequency |
|-----------|----------------|------------------|
| Kubernetes | Monthly | Quarterly |
| Application dependencies | Monthly | As needed |
| Base images | Weekly | Monthly |
| Infrastructure (DB, Redis, etc.) | Monthly | Quarterly |
| Security patches | Daily | Immediately (critical) |
| TLS certificates | Daily | Auto (cert-manager) |

---

## Appendix

### A. Quick Reference

#### Essential Commands

```bash
# Status checks
kubectl get all -n essp-prod
kubectl get pods -n essp-prod -o wide
kubectl top pods -n essp-prod

# Logs
kubectl logs -n essp-prod -l app=ims-api --tail=100 -f

# Scale
kubectl scale deployment ims-api -n essp-prod --replicas=5

# Restart
kubectl rollout restart deployment/ims-api -n essp-prod

# Rollback
helm rollback essp -n essp-prod
kubectl rollout undo deployment/ims-api -n essp-prod

# Port-forward
kubectl port-forward -n essp-prod svc/ims-api 8080:8080

# Exec
kubectl exec -it <pod-name> -n essp-prod -- /bin/sh
```

### B. Contact Information

| Role | Contact | Escalation |
|------|---------|------------|
| Platform Team | platform@essp.example.com | Slack: #essp-platform |
| Database Admin | dba@essp.example.com | PagerDuty: essp-dba |
| Security Team | security@essp.example.com | Slack: #security |
| On-Call Engineer | oncall@essp.example.com | PagerDuty: essp-oncall |

### C. External Resources

- **Kubernetes Documentation**: https://kubernetes.io/docs/
- **Helm Documentation**: https://helm.sh/docs/
- **Prometheus**: https://prometheus.io/docs/
- **Grafana**: https://grafana.com/docs/
- **PostgreSQL**: https://www.postgresql.org/docs/
- **ESSP Repository**: <repository-url>

### D. Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2025-12-12 | Initial deployment runbook |

---

**Document Owner**: Platform Team
**Last Reviewed**: 2025-12-12
**Next Review**: 2026-01-12
