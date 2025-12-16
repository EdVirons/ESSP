# Changelog

All notable changes to the ESSP Helm chart will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-12-12

### Added
- Initial release of ESSP Helm chart
- Support for all ESSP microservices:
  - IMS API service with HPA and PDB
  - SSOT School service
  - SSOT Devices service
  - SSOT Parts service
  - Sync Worker service
- Environment-specific value files (dev, staging, prod)
- Comprehensive configuration options for:
  - PostgreSQL database
  - Redis cache
  - NATS messaging
  - MinIO object storage
- Ingress configuration with TLS support
- Network policies for enhanced security
- ConfigMap and Secret management
- Service Account creation
- Health check probes for all services
- Resource limits and requests
- Horizontal Pod Autoscaling for IMS API
- Pod Disruption Budget for high availability
- Security contexts and pod security policies
- Helper templates for common functions
- Comprehensive README with installation and configuration guide
- Example installation scripts for different environments
- Custom values example file

### Security
- Non-root container execution
- Read-only root filesystem
- Dropped all capabilities
- Security contexts enforced
- Network policies support

## [Unreleased]

### Planned
- Support for external secrets operator
- Database migration job
- Backup and restore procedures
- Monitoring and alerting integration
- Service mesh support (Istio/Linkerd)
- Additional environment presets
- Chart dependencies for PostgreSQL, Redis, NATS, MinIO
- Init containers for database migration
- Custom resource definitions (CRDs) if needed
