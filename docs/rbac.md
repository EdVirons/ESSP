# RBAC & Apps

## Apps
- Field Tech App
- Demo/Sales App
- Support Desk App
- School IT Visibility App
- Supplier Portal
- Contractor Portal
- County Service Shop (Warehouse) App

## Suggested roles (Keycloak)
- `ssp_admin` — tenant super admin
- `ssp_support_agent` — tickets/dispatch
- `ssp_field_tech` — work orders + deliverables
- `ssp_lead_tech` — scheduling + approvals request
- `ssp_demo_team` — demos, surveys, pipeline
- `ssp_sales` — pipeline visibility
- `ssp_school_contact` — create incidents, approve sign-offs
- `ssp_supplier` — parts catalog + fulfillment visibility
- `ssp_contractor` — work packages, deliverables submission
- `ssp_warehouse_manager` — inventory, BOM pick/issue

## Authorization pattern
- Role check + `tenant_id`
- School-scoped endpoints also check:
  - user has access to `school_id` (claim or mapping)
