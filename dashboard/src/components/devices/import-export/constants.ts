export const exportFields = [
  { key: 'serial', label: 'Serial Number', default: true },
  { key: 'assetTag', label: 'Asset Tag', default: true },
  { key: 'make', label: 'Make', default: true },
  { key: 'model', label: 'Model', default: true },
  { key: 'category', label: 'Category', default: true },
  { key: 'schoolId', label: 'School ID', default: true },
  { key: 'schoolName', label: 'School Name', default: true },
  { key: 'lifecycle', label: 'Status', default: true },
  { key: 'enrolled', label: 'Enrollment', default: true },
  { key: 'assignedTo', label: 'Assigned To', default: false },
  { key: 'purchaseDate', label: 'Purchase Date', default: false },
  { key: 'warrantyExpiry', label: 'Warranty Expiry', default: false },
  { key: 'lastSeen', label: 'Last Seen', default: false },
  { key: 'notes', label: 'Notes', default: false },
  { key: 'createdAt', label: 'Created At', default: false },
];

export type ExportScope = 'all' | 'filtered' | 'selected';
