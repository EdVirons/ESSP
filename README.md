# EdVirons School Services Platform (SSP)

A Go-based, multi-tenant platform to manage **school edtech services** end-to-end:
- **SSOT** (School, Devices, Parts/PUK) as authoritative reference data
- **IMS** operational workflows: incidents → tickets → work orders → deliverables → approvals
- **Rollout lifecycle**: demo → survey → install (low voltage) → integrate → commission → ops
- **Warehouses/service-shops** per county: inventory, BOM, allocations, fulfillment

This repo is designed as a **monorepo** with multiple services + shared packages.

## Services
- `services/ims-api` — incidents, tickets, work orders, rollouts, approvals, deliverables
- `services/ssot-school` — School SSOT (imports/sync from registry), counties/sub-counties, school contacts
- `services/ssot-devices` — Devices SSOT (inventory, assignments, lifecycle), model catalog
- `services/ssot-parts` — Parts/PUK SSOT (catalog, compat, vendor SKUs)
- `services/sync-worker` — pulls SSOT snapshots into IMS read models / caches (optional)

## Infra (local dev)
- PostgreSQL (data)
- Valkey (cache)
- MinIO (attachments/evidence)
- NATS (events) — optional but included for clean decoupling

## Quick start
```bash
cd deployments/docker
docker compose up -d --build
```

Run migrations per service:
```bash
make -C services/ims-api migrate-up
make -C services/ssot-school migrate-up
make -C services/ssot-devices migrate-up
make -C services/ssot-parts migrate-up
```

## Docs
See `docs/`:
- `docs/architecture.md` — system overview, boundaries, data ownership
- `docs/rbac.md` — roles, apps, permissions
- `docs/api.md` — routing conventions + headers
- `docs/ssot.md` — SSOT design rules and sync strategies
