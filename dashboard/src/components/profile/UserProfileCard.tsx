import { Mail, Building2, Shield, Clock, CheckCircle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { UserAvatar } from './UserAvatar';
import type { SSOUserProfile } from '@/lib/api';

interface UserProfileCardProps {
  profile: SSOUserProfile;
  compact?: boolean;
}

// Map role names to display labels
const roleLabels: Record<string, { label: string; variant: 'default' | 'secondary' | 'outline'; color: string }> = {
  ssp_admin: { label: 'Administrator', variant: 'default', color: 'from-cyan-500 to-teal-500' },
  ssp_operator: { label: 'Operator', variant: 'secondary', color: 'from-blue-500 to-indigo-500' },
  ssp_viewer: { label: 'Viewer', variant: 'outline', color: 'from-gray-400 to-gray-500' },
};

function formatDate(dateString?: string): string {
  if (!dateString) return 'Never';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function UserProfileCard({ profile, compact = false }: UserProfileCardProps) {
  if (compact) {
    return (
      <div className="flex items-center gap-3 p-3">
        <UserAvatar
          src={profile.avatarUrl}
          fallback={profile.displayName}
          size="md"
        />
        <div className="min-w-0 flex-1">
          <p className="truncate text-sm font-medium text-gray-900">
            {profile.displayName}
          </p>
          <p className="truncate text-xs text-gray-500">
            {profile.email || profile.username}
          </p>
        </div>
      </div>
    );
  }

  return (
    <Card className="overflow-hidden border-0 shadow-md">
      {/* Gradient banner */}
      <div className="h-20 bg-gradient-to-r from-cyan-500 via-teal-500 to-cyan-600 relative">
        <div className="absolute inset-0 opacity-20">
          <div className="absolute top-2 right-8 w-16 h-16 rounded-full bg-white blur-2xl" />
          <div className="absolute bottom-0 left-12 w-24 h-12 rounded-full bg-teal-300 blur-xl" />
        </div>
      </div>
      <CardHeader className="pb-4 -mt-10 relative">
        <div className="flex items-end gap-4">
          <div className="rounded-full bg-white p-1 shadow-lg">
            <UserAvatar
              src={profile.avatarUrl}
              fallback={profile.displayName}
              size="xl"
            />
          </div>
          <div className="flex-1 space-y-1 pb-1">
            <CardTitle className="text-xl">{profile.displayName}</CardTitle>
            <p className="text-sm text-gray-500">@{profile.username}</p>
          </div>
        </div>
        <div className="flex flex-wrap gap-1.5 pt-3">
          {profile.roles.map((role) => {
            const roleInfo = roleLabels[role] || { label: role, variant: 'outline' as const, color: 'from-gray-400 to-gray-500' };
            return (
              <Badge key={role} variant={roleInfo.variant} className="shadow-sm">
                <Shield className="mr-1 h-3 w-3" />
                {roleInfo.label}
              </Badge>
            );
          })}
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {/* Email */}
        {profile.email && (
          <div className="flex items-center gap-3 text-sm rounded-lg bg-gray-50 p-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-cyan-100">
              <Mail className="h-4 w-4 text-cyan-600" />
            </div>
            <span className="text-gray-700 flex-1">{profile.email}</span>
            {profile.emailVerified && (
              <span title="Verified" className="flex items-center gap-1 text-xs text-emerald-600">
                <CheckCircle className="h-4 w-4" />
                Verified
              </span>
            )}
          </div>
        )}

        {/* Organization */}
        {profile.organization && (
          <div className="flex items-center gap-3 text-sm rounded-lg bg-gray-50 p-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-100">
              <Building2 className="h-4 w-4 text-purple-600" />
            </div>
            <div>
              <span className="text-gray-900 font-medium">{profile.organization.displayName || profile.organization.name}</span>
              {profile.organization.type && (
                <span className="ml-2 text-gray-500 text-xs">({profile.organization.type})</span>
              )}
            </div>
          </div>
        )}

        {/* Last Login */}
        {profile.lastLoginAt && (
          <div className="flex items-center gap-3 text-sm rounded-lg bg-gray-50 p-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-100">
              <Clock className="h-4 w-4 text-amber-600" />
            </div>
            <span className="text-gray-600">
              Last login: <span className="text-gray-900">{formatDate(profile.lastLoginAt)}</span>
            </span>
          </div>
        )}

        {/* SSO Provider info */}
        {profile.ssoProvider && (
          <div className="mt-2 rounded-lg bg-gradient-to-r from-cyan-50 to-teal-50 border border-cyan-100 p-3">
            <p className="text-xs text-cyan-700">
              Signed in via <span className="font-semibold">{profile.ssoProvider}</span>
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
