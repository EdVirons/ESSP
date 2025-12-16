# API conventions

## Headers
- `X-Tenant-Id` (required)
- `X-School-Id` (required on school-scoped IMS endpoints)
- `Authorization: Bearer <JWT>` (recommended)

## Response format
- JSON
- Lists return: `{ "items": [...], "nextCursor": "..." }` when paginated.

## Cursors
- Opaque cursor: base64 encoded `created_at|id`


## SSOT exports
Each SSOT service supports:
- `GET /v1/export`
- `POST /v1/import`

Header: `X-Tenant-Id` is required.

On import, SSOT publishes to NATS:
- `ssot.<kind>.snapshot`
- `ssot.<kind>.changed`


## IMS SSOT cache
IMS stores last-seen SSOT exports in Postgres table `ims_ssot_snapshots` (JSONB) keyed by `(tenant_id, kind)`.
This supports fast lookups inside IMS even if SSOT services are temporarily unavailable.


## IMS SSOT lookup endpoints (debug)
These read from IMS cached snapshots table `ims_ssot_snapshots`.

- `GET /v1/ssot/school/{id}`
- `GET /v1/ssot/school/{id}/primary-contact`
- `GET /v1/ssot/device/{id}`
- `GET /v1/ssot/device/serial/{serial}`
- `GET /v1/ssot/part/{id}`
- `GET /v1/ssot/part/puk/{puk}`

Header: `X-Tenant-Id` required.


## BOM suggestions
- `GET /v1/work-orders/{id}/bom/suggest?q=&limit=`

Returns parts compatible with the work order's device model (from SSOT snapshots).
Header: `X-Tenant-Id`, `X-School-Id` required.


## BOM compatibility enforcement
By default, `POST /v1/work-orders/{id}/bom` enforces part compatibility vs the work order device model (from SSOT snapshots).

- Disable enforcement: set env `BOM_ENFORCE_COMPATIBILITY=false`
- One-off override (when enforcement enabled): `?allowIncompatible=true`
