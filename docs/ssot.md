# SSOT design rules

## What belongs in SSOT
- Stable reference datasets
- Canonical identifiers used by other services
- Import/sync from official registries

## What does NOT belong in SSOT
- Operational data (incidents, work orders)
- Evidence artifacts (photos, pdfs)
- Approval decisions, schedules

## Sync strategies
1) Event-driven (recommended): SSOT publishes events, IMS subscribes
2) Snapshot pull: `sync-worker` pulls SSOT exports on a schedule
3) Hybrid: events for changes + nightly reconcile snapshots

## ID rules
- SSOT generates IDs, never IMS.
- IMS stores only SSOT IDs + human-readable denormalized fields (optional).
