import { useEffect } from 'react';
import { Loader2, User, ArrowLeft, UserCircle, Sparkles } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { UserProfileCard, UserPermissions, UserPreferencesCard } from '@/components/profile';

export function Profile() {
  const { user, profile, isLoadingProfile, fetchProfile } = useAuth();

  // Fetch profile data on mount
  useEffect(() => {
    if (user && !profile) {
      fetchProfile();
    }
  }, [user, profile, fetchProfile]);

  if (isLoadingProfile) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-cyan-600" />
      </div>
    );
  }

  // Show a simple view if profile hasn't loaded yet
  if (!profile && user) {
    return (
      <div className="space-y-6">
        {/* Header */}
        <div className="rounded-xl bg-gradient-to-r from-cyan-600 via-teal-600 to-cyan-700 p-6 text-white shadow-lg">
          <div className="flex items-center gap-4">
            <Link to="/overview">
              <Button variant="ghost" size="sm" className="text-white/80 hover:text-white hover:bg-white/10">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back
              </Button>
            </Link>
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
                <UserCircle className="h-6 w-6" />
              </div>
              <div>
                <h1 className="text-2xl font-bold">User Profile</h1>
                <p className="text-cyan-100">Your account information</p>
              </div>
            </div>
          </div>
        </div>

        <div className="rounded-xl border-0 bg-white p-6 shadow-md">
          <div className="flex items-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-gradient-to-br from-cyan-100 to-teal-100">
              <User className="h-8 w-8 text-cyan-600" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-gray-900">
                {user.displayName || user.username}
              </h2>
              <p className="text-sm text-gray-500">@{user.username}</p>
              {user.email && (
                <p className="text-sm text-gray-500">{user.email}</p>
              )}
            </div>
          </div>
          <div className="mt-4 flex flex-wrap gap-2">
            {user.roles?.map((role) => (
              <span
                key={role}
                className="rounded-full bg-gradient-to-r from-cyan-100 to-teal-100 px-3 py-1 text-xs font-medium text-cyan-800"
              >
                {role}
              </span>
            ))}
          </div>
        </div>

        <div className="flex justify-center">
          <Button onClick={() => fetchProfile()}>
            Load Full Profile
          </Button>
        </div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="flex h-64 items-center justify-center">
        <p className="text-gray-500">Unable to load profile</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="rounded-xl bg-gradient-to-r from-cyan-600 via-teal-600 to-cyan-700 p-6 text-white shadow-lg">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link to="/overview">
              <Button variant="ghost" size="sm" className="text-white/80 hover:text-white hover:bg-white/10">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back
              </Button>
            </Link>
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/20 backdrop-blur">
                <UserCircle className="h-6 w-6" />
              </div>
              <div>
                <h1 className="text-2xl font-bold">User Profile</h1>
                <p className="text-cyan-100">Your account information and preferences</p>
              </div>
            </div>
          </div>
          <div className="hidden md:flex items-center gap-2 rounded-lg bg-white/10 px-3 py-1.5 backdrop-blur">
            <Sparkles className="h-4 w-4 text-cyan-200" />
            <span className="text-sm text-cyan-100">ESSP Account</span>
          </div>
        </div>
      </div>

      {/* Main content grid */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Profile card */}
        <UserProfileCard profile={profile} />

        {/* Preferences card */}
        <UserPreferencesCard
          preferences={profile.preferences}
          onEdit={() => {
            // TODO: Open preferences edit modal
            console.log('Edit preferences');
          }}
        />

        {/* Permissions card - full width */}
        <div className="lg:col-span-2">
          <UserPermissions permissions={profile.permissions || []} />
        </div>
      </div>

      {/* Session information */}
      <div className="rounded-xl border border-cyan-100 bg-gradient-to-r from-cyan-50 to-teal-50 p-4 shadow-sm">
        <h3 className="mb-3 text-sm font-semibold text-cyan-800">Session Information</h3>
        <div className="grid gap-4 text-sm sm:grid-cols-2 lg:grid-cols-4">
          <div className="rounded-lg bg-white/60 p-2 backdrop-blur">
            <span className="text-cyan-600 text-xs font-medium">Tenant ID</span>
            <p className="font-mono text-gray-900 text-xs mt-0.5 truncate">{profile.tenantId}</p>
          </div>
          <div className="rounded-lg bg-white/60 p-2 backdrop-blur">
            <span className="text-cyan-600 text-xs font-medium">User ID</span>
            <p className="font-mono text-gray-900 text-xs mt-0.5 truncate">{profile.id}</p>
          </div>
          {profile.ssoSubject && (
            <div className="rounded-lg bg-white/60 p-2 backdrop-blur">
              <span className="text-cyan-600 text-xs font-medium">SSO Subject</span>
              <p className="font-mono text-gray-900 text-xs mt-0.5 truncate">{profile.ssoSubject}</p>
            </div>
          )}
          {profile.createdAt && (
            <div className="rounded-lg bg-white/60 p-2 backdrop-blur">
              <span className="text-cyan-600 text-xs font-medium">Account Created</span>
              <p className="text-gray-900 text-xs mt-0.5">
                {new Date(profile.createdAt).toLocaleDateString()}
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
