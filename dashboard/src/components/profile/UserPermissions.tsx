import { Shield, Check, Eye, Plus, Pencil, Trash2 } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';

interface UserPermissionsProps {
  permissions: string[];
}

// Group permissions by resource
function groupPermissions(permissions: string[]): Record<string, string[]> {
  const groups: Record<string, string[]> = {};

  for (const perm of permissions) {
    const [resource, action] = perm.split(':');
    if (!groups[resource]) {
      groups[resource] = [];
    }
    groups[resource].push(action);
  }

  return groups;
}

// Resource display names
const resourceLabels: Record<string, string> = {
  incident: 'Incidents',
  workorder: 'Work Orders',
  program: 'Programs',
  school: 'Schools',
  device: 'Devices',
  serviceshop: 'Service Shops',
  inventory: 'Inventory',
  audit: 'Audit Logs',
  settings: 'Settings',
  user: 'Users',
};

// Action icons/colors
const actionStyles: Record<string, { color: string; bg: string; icon: typeof Eye }> = {
  read: { color: 'text-cyan-700', bg: 'bg-cyan-100', icon: Eye },
  create: { color: 'text-emerald-700', bg: 'bg-emerald-100', icon: Plus },
  update: { color: 'text-amber-700', bg: 'bg-amber-100', icon: Pencil },
  delete: { color: 'text-red-700', bg: 'bg-red-100', icon: Trash2 },
};

export function UserPermissions({ permissions }: UserPermissionsProps) {
  const grouped = groupPermissions(permissions);
  const resources = Object.keys(grouped).sort();

  if (permissions.length === 0) {
    return (
      <Card className="border-0 shadow-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-500 to-teal-600">
              <Shield className="h-4 w-4 text-white" />
            </div>
            Permissions
          </CardTitle>
          <CardDescription>
            No specific permissions assigned
          </CardDescription>
        </CardHeader>
      </Card>
    );
  }

  return (
    <Card className="border-0 shadow-md">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-500 to-teal-600">
            <Shield className="h-4 w-4 text-white" />
          </div>
          Permissions
        </CardTitle>
        <CardDescription>
          Access rights based on your role
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {resources.map((resource) => (
            <div key={resource} className="rounded-xl border border-gray-100 bg-gradient-to-br from-gray-50 to-white p-4 hover:shadow-sm transition-shadow">
              <h4 className="mb-3 text-sm font-semibold text-gray-900">
                {resourceLabels[resource] || resource}
              </h4>
              <div className="flex flex-wrap gap-1.5">
                {grouped[resource].map((action) => {
                  const style = actionStyles[action] || { color: 'text-gray-600', bg: 'bg-gray-100', icon: Check };
                  const Icon = style.icon;
                  return (
                    <span
                      key={action}
                      className={`inline-flex items-center gap-1 rounded-lg px-2 py-1 text-xs font-medium ${style.bg} ${style.color}`}
                    >
                      <Icon className="h-3 w-3" />
                      {action}
                    </span>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
