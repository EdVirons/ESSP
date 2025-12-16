# Architecture

## Principle
- **SSOT services own truth** for stable datasets (schools, devices, parts).
- **IMS owns operations** (workflows, evidence, approvals, scheduling).
- **SSOT IDs are referenced**, not duplicated, in IMS (e.g., `school_id`, `device_id`, `part_id`).

## Service boundaries
### SSOT School
Owns:
- Counties/sub-counties/wards (Kenya admin units)
- School registry fields + immutable identifiers
- Default school contacts (Point-of-Contact records can be mirrored into IMS)

### SSOT Devices
Owns:
- Device inventory, lifecycle state, assignments, device model catalog

### SSOT Parts (PUK)
Owns:
- Part catalog, compatibility, vendor SKUs, pricing bands (authoritative)

### IMS
Owns:
- Incidents / tickets
- Work orders: scheduling, deliverables (MinIO evidence), approvals
- BOM allocations to work orders; fulfillment from county service shops
- Rollout programs & phases, with **phase gates**

## Data flow
- SSOT services publish **events** (`school.updated`, `device.registered`, `part.updated`) to NATS.
- IMS subscribes (or `sync-worker` does) and maintains:
  - read-optimized lookup tables (optional)
  - cache in Valkey for hot paths
- Evidence and attachments stored in MinIO; metadata in Postgres.

## Multi-tenancy
- Every table includes `tenant_id`.
- Requests require `X-Tenant-Id`. IMS routes additionally require `X-School-Id` for school-scoped ops.

## Security
- Use Keycloak/OIDC in front of services (or API Gateway). Services validate JWTs.
- Authorization via role-based access control (RBAC) + school scoping.


## SSOT consumption (IMS)
IMS subscribes to NATS subjects:
- `ssot.school.changed`
- `ssot.devices.changed`
- `ssot.parts.changed`

On change, IMS fetches the latest SSOT export (`GET /v1/export`) from the relevant SSOT service
and upserts it into `ims_ssot_snapshots` in the IMS database (JSONB).
