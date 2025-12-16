// Report filters
export interface ReportFilters {
  dateFrom?: string;
  dateTo?: string;
  status?: string[];
  schoolId?: string;
  countyCode?: string;
  category?: string;
  sortBy?: string;
  sortDir?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
}

// Pagination
export interface ReportPagination {
  total: number;
  offset: number;
  limit: number;
}

// Work Orders Report
export interface WorkOrderReportItem {
  id: string;
  incidentId: string;
  status: string;
  taskType: string;
  schoolName: string;
  deviceCategory: string;
  assignedTo: string;
  costCents: number;
  reworkCount: number;
  createdAt: string;
  completedAt?: string;
  durationHours?: number;
}

export interface WorkOrderReportSummary {
  total: number;
  byStatus: Record<string, number>;
  avgCompletionHours: number;
  totalCostCents: number;
  reworkRate: number;
}

export interface WorkOrderReportResponse {
  items: WorkOrderReportItem[];
  summary: WorkOrderReportSummary;
  pagination: ReportPagination;
}

// Incidents Report
export interface IncidentReportItem {
  id: string;
  title: string;
  status: string;
  severity: string;
  category: string;
  schoolName: string;
  slaBreached: boolean;
  createdAt: string;
  resolvedAt?: string;
  resolutionHours?: number;
}

export interface IncidentReportSummary {
  total: number;
  byStatus: Record<string, number>;
  bySeverity: Record<string, number>;
  slaBreachedCount: number;
  slaComplianceRate: number;
  avgResolutionHours: number;
}

export interface IncidentReportResponse {
  items: IncidentReportItem[];
  summary: IncidentReportSummary;
  pagination: ReportPagination;
}

// Inventory Report
export interface InventoryReportItem {
  partId: string;
  partSku: string;
  partName: string;
  category: string;
  serviceShopName: string;
  qtyAvailable: number;
  qtyReserved: number;
  reorderThreshold: number;
  isLowStock: boolean;
}

export interface InventoryReportSummary {
  totalParts: number;
  lowStockCount: number;
  totalQtyAvailable: number;
  byCategory: Record<string, number>;
}

export interface InventoryReportResponse {
  items: InventoryReportItem[];
  summary: InventoryReportSummary;
  pagination: ReportPagination;
}

// Schools Report
export interface SchoolReportItem {
  schoolId: string;
  schoolName: string;
  countyName: string;
  deviceCount: number;
  incidentCount: number;
  workOrderCount: number;
}

export interface SchoolReportSummary {
  totalSchools: number;
  totalDevices: number;
  byCounty: Record<string, number>;
}

export interface SchoolReportResponse {
  items: SchoolReportItem[];
  summary: SchoolReportSummary;
  pagination: ReportPagination;
}

// Executive Dashboard
export interface ExecutiveDashboard {
  workOrders: {
    total: number;
    completed: number;
    inProgress: number;
    completionRate: number;
    avgCompletionDays: number;
  };
  incidents: {
    total: number;
    open: number;
    resolved: number;
    slaCompliance: number;
    critical: number;
  };
  inventory: {
    totalParts: number;
    lowStock: number;
    outOfStock: number;
  };
  schools: {
    totalSchools: number;
    totalDevices: number;
    activeProjects: number;
  };
}
