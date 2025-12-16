# ESSP Management Dashboard

A comprehensive web-based administration interface for the ESSP (EdVirons School Services Platform) microservices platform.

## Overview

The ESSP Dashboard provides operations teams with:

- **Service Health Monitoring** - Real-time status of all microservices
- **Key Metrics** - Incidents, work orders, and program statistics
- **Data Management** - CRUD operations for core business entities
- **Audit Trail** - Complete activity and change history
- **Configuration** - Environment and feature settings

## Technology Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS v4** for styling
- **TanStack Query** for data fetching and caching
- **React Router v6** for navigation
- **Lucide React** for icons

## Getting Started

### Prerequisites

- Node.js 18+ or 20+
- npm or pnpm

### Development

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start development server:
   ```bash
   npm run dev
   ```

3. Open http://localhost:5173 in your browser

### Building

```bash
npm run build
```

The production build will be in the `dist/` directory.

### Linting

```bash
npm run lint
```

## Project Structure

```
src/
  api/            # API client and hooks
  components/
    layout/       # Layout components (Header, Sidebar, etc.)
    overview/     # Dashboard overview components
    incidents/    # Incident management components
    work-orders/  # Work order components
    ui/           # Base UI components
  hooks/          # Custom React hooks
  lib/            # Utility functions and constants
  pages/          # Page components
  types/          # TypeScript type definitions
```

## Available Routes

| Route | Description |
|-------|-------------|
| `/overview` | Dashboard home with metrics and health |
| `/incidents` | Incident list and management |
| `/incidents/:id` | Incident detail view |
| `/work-orders` | Work order list and management |
| `/work-orders/:id` | Work order detail view |
| `/programs` | Program list and management |
| `/service-shops` | Service shop management |
| `/audit-logs` | System activity logs |
| `/settings` | User and system settings |

## API Integration

The dashboard connects to the IMS API backend. Configure the API endpoint in development mode via the Vite proxy configuration in `vite.config.ts`.

Default proxy settings:
- `/api/*` routes to `http://localhost:8080/v1/*`
- `/admin/*` routes to `http://localhost:8080/admin/*`

## Backend Admin Endpoints

The dashboard requires the following admin endpoints (in ims-api):

- `GET /admin/v1/health/services` - Aggregated service health
- `GET /admin/v1/metrics/summary` - Dashboard metrics summary
- `GET /admin/v1/activity` - Recent activity feed

## Documentation

- [Architecture Design](../docs/dashboard/ARCHITECTURE.md)
- [Implementation Plan](../docs/dashboard/IMPLEMENTATION_PLAN.md)

## Development Notes

### Adding New Pages

1. Create the page component in `src/pages/`
2. Add the route in `src/App.tsx`
3. Add navigation item in `src/components/layout/Sidebar.tsx`

### Adding New API Endpoints

1. Add types in `src/types/index.ts`
2. Create API hooks in `src/api/`
3. Use TanStack Query patterns for caching

### UI Components

Base UI components are in `src/components/ui/` and follow the shadcn/ui patterns.

## License

Copyright (c) 2024 EdVirons. All rights reserved.
