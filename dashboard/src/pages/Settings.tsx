import * as React from 'react';
import {
  Settings as SettingsIcon,
  User,
  Bell,
  Palette,
  Shield,
  Database,
  RefreshCw,
  Eye,
  EyeOff,
  Save,
  Moon,
  Sun,
  Monitor,
  Wifi,
  Clock,
  HardDrive,
  Cloud,
  Key,
  LogOut,
  Trash2,
  Download,
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { toast } from 'sonner';

// Toggle Switch Component
function ToggleSwitch({
  checked,
  onChange,
  disabled,
  label,
}: {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
  label?: string;
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      aria-label={label}
      onClick={() => !disabled && onChange(!checked)}
      disabled={disabled}
      className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
        checked ? 'bg-blue-600' : 'bg-gray-200'
      } ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
    >
      <span
        className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
          checked ? 'translate-x-6' : 'translate-x-1'
        }`}
      />
    </button>
  );
}

// Settings Row Component
function SettingsRow({
  title,
  description,
  children,
}: {
  title: string;
  description?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between py-3 border-b border-gray-100 last:border-0">
      <div>
        <p className="text-sm font-medium text-gray-900">{title}</p>
        {description && <p className="text-xs text-gray-500">{description}</p>}
      </div>
      <div>{children}</div>
    </div>
  );
}

// Status Indicator
function StatusIndicator({ status, label }: { status: 'online' | 'offline' | 'warning'; label: string }) {
  const colors = {
    online: 'bg-green-500',
    offline: 'bg-red-500',
    warning: 'bg-yellow-500',
  };
  return (
    <div className="flex items-center gap-2">
      <span className={`h-2 w-2 rounded-full ${colors[status]}`} />
      <span className="text-sm text-gray-600">{label}</span>
    </div>
  );
}

export function Settings() {
  // Active tab state
  const [activeTab, setActiveTab] = React.useState('profile');

  // Profile state
  const [profile, setProfile] = React.useState({
    name: 'Admin User',
    email: 'admin@essp.local',
    phone: '+254 700 000 000',
  });
  const [isEditingProfile, setIsEditingProfile] = React.useState(false);
  const [showPassword, setShowPassword] = React.useState(false);
  const [passwordForm, setPasswordForm] = React.useState({
    current: '',
    new: '',
    confirm: '',
  });

  // Notification settings
  const [notifications, setNotifications] = React.useState({
    emailAlerts: true,
    browserNotifications: false,
    workOrderUpdates: true,
    incidentAlerts: true,
    syncNotifications: false,
    weeklyReports: true,
  });

  // Appearance settings
  const [appearance, setAppearance] = React.useState({
    theme: 'light' as 'light' | 'dark' | 'system',
    sidebarCollapsed: false,
    compactMode: false,
    showAvatars: true,
  });

  // System status (mock data - would come from API)
  const [systemStatus] = React.useState({
    api: 'online' as const,
    database: 'online' as const,
    redis: 'online' as const,
    nats: 'online' as const,
    minio: 'online' as const,
    ssot: 'online' as const,
  });

  // SSOT settings
  const [ssotSettings, setSsotSettings] = React.useState({
    autoSync: true,
    syncInterval: 30,
    syncOnStartup: true,
    conflictResolution: 'ssot_wins' as 'ssot_wins' | 'local_wins' | 'manual',
  });

  // Last sync time
  const [lastSync] = React.useState(new Date().toISOString());

  const handleSaveProfile = () => {
    toast.success('Profile updated successfully');
    setIsEditingProfile(false);
  };

  const handleChangePassword = () => {
    if (passwordForm.new !== passwordForm.confirm) {
      toast.error('Passwords do not match');
      return;
    }
    if (passwordForm.new.length < 8) {
      toast.error('Password must be at least 8 characters');
      return;
    }
    toast.success('Password changed successfully');
    setPasswordForm({ current: '', new: '', confirm: '' });
  };

  const handleClearCache = () => {
    localStorage.clear();
    sessionStorage.clear();
    toast.success('Cache cleared successfully');
  };

  const handleExportData = () => {
    toast.success('Data export started. You will receive an email when ready.');
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
        <p className="text-sm text-gray-500">Manage your account and application preferences</p>
      </div>

      {/* Settings Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="bg-white border border-gray-200 p-1">
          <TabsTrigger value="profile" className="gap-2">
            <User className="h-4 w-4" />
            Profile
          </TabsTrigger>
          <TabsTrigger value="notifications" className="gap-2">
            <Bell className="h-4 w-4" />
            Notifications
          </TabsTrigger>
          <TabsTrigger value="appearance" className="gap-2">
            <Palette className="h-4 w-4" />
            Appearance
          </TabsTrigger>
          <TabsTrigger value="system" className="gap-2">
            <SettingsIcon className="h-4 w-4" />
            System
          </TabsTrigger>
          <TabsTrigger value="security" className="gap-2">
            <Shield className="h-4 w-4" />
            Security
          </TabsTrigger>
        </TabsList>

        {/* Profile Tab */}
        <TabsContent value="profile" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  Personal Information
                </CardTitle>
                <CardDescription>Update your personal details</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
                    {isEditingProfile ? (
                      <Input
                        value={profile.name}
                        onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                      />
                    ) : (
                      <p className="text-sm text-gray-900 py-2">{profile.name}</p>
                    )}
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
                    {isEditingProfile ? (
                      <Input
                        type="email"
                        value={profile.email}
                        onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                      />
                    ) : (
                      <p className="text-sm text-gray-900 py-2">{profile.email}</p>
                    )}
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
                    {isEditingProfile ? (
                      <Input
                        value={profile.phone}
                        onChange={(e) => setProfile({ ...profile, phone: e.target.value })}
                      />
                    ) : (
                      <p className="text-sm text-gray-900 py-2">{profile.phone}</p>
                    )}
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Role</label>
                    <Badge className="bg-blue-100 text-blue-800">SSP Admin</Badge>
                  </div>
                  <div className="pt-2">
                    {isEditingProfile ? (
                      <div className="flex gap-2">
                        <Button onClick={handleSaveProfile}>
                          <Save className="h-4 w-4 mr-2" />
                          Save Changes
                        </Button>
                        <Button variant="outline" onClick={() => setIsEditingProfile(false)}>
                          Cancel
                        </Button>
                      </div>
                    ) : (
                      <Button variant="outline" onClick={() => setIsEditingProfile(true)}>
                        Edit Profile
                      </Button>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Key className="h-5 w-5" />
                  Change Password
                </CardTitle>
                <CardDescription>Update your account password</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Current Password</label>
                    <div className="relative">
                      <Input
                        type={showPassword ? 'text' : 'password'}
                        value={passwordForm.current}
                        onChange={(e) => setPasswordForm({ ...passwordForm, current: e.target.value })}
                        placeholder="Enter current password"
                      />
                      <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                        aria-label={showPassword ? 'Hide password' : 'Show password'}
                      >
                        {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </button>
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">New Password</label>
                    <Input
                      type={showPassword ? 'text' : 'password'}
                      value={passwordForm.new}
                      onChange={(e) => setPasswordForm({ ...passwordForm, new: e.target.value })}
                      placeholder="Enter new password"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Confirm New Password</label>
                    <Input
                      type={showPassword ? 'text' : 'password'}
                      value={passwordForm.confirm}
                      onChange={(e) => setPasswordForm({ ...passwordForm, confirm: e.target.value })}
                      placeholder="Confirm new password"
                    />
                  </div>
                  <Button
                    onClick={handleChangePassword}
                    disabled={!passwordForm.current || !passwordForm.new || !passwordForm.confirm}
                  >
                    Update Password
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Notifications Tab */}
        <TabsContent value="notifications" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Bell className="h-5 w-5" />
                  Notification Channels
                </CardTitle>
                <CardDescription>Choose how you want to be notified</CardDescription>
              </CardHeader>
              <CardContent>
                <SettingsRow title="Email Notifications" description="Receive alerts via email">
                  <ToggleSwitch
                    checked={notifications.emailAlerts}
                    onChange={(checked) => setNotifications({ ...notifications, emailAlerts: checked })}
                    label="Toggle email notifications"
                  />
                </SettingsRow>
                <SettingsRow title="Browser Notifications" description="Show desktop push notifications">
                  <ToggleSwitch
                    checked={notifications.browserNotifications}
                    onChange={(checked) => setNotifications({ ...notifications, browserNotifications: checked })}
                    label="Toggle browser notifications"
                  />
                </SettingsRow>
                <SettingsRow title="Weekly Reports" description="Receive weekly summary emails">
                  <ToggleSwitch
                    checked={notifications.weeklyReports}
                    onChange={(checked) => setNotifications({ ...notifications, weeklyReports: checked })}
                    label="Toggle weekly reports"
                  />
                </SettingsRow>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <RefreshCw className="h-5 w-5" />
                  Notification Types
                </CardTitle>
                <CardDescription>Select which events trigger notifications</CardDescription>
              </CardHeader>
              <CardContent>
                <SettingsRow title="Work Order Updates" description="When work orders are created or updated">
                  <ToggleSwitch
                    checked={notifications.workOrderUpdates}
                    onChange={(checked) => setNotifications({ ...notifications, workOrderUpdates: checked })}
                    label="Toggle work order updates"
                  />
                </SettingsRow>
                <SettingsRow title="Incident Alerts" description="When new incidents are reported">
                  <ToggleSwitch
                    checked={notifications.incidentAlerts}
                    onChange={(checked) => setNotifications({ ...notifications, incidentAlerts: checked })}
                    label="Toggle incident alerts"
                  />
                </SettingsRow>
                <SettingsRow title="Sync Notifications" description="When SSOT sync completes">
                  <ToggleSwitch
                    checked={notifications.syncNotifications}
                    onChange={(checked) => setNotifications({ ...notifications, syncNotifications: checked })}
                    label="Toggle sync notifications"
                  />
                </SettingsRow>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Appearance Tab */}
        <TabsContent value="appearance" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Palette className="h-5 w-5" />
                  Theme
                </CardTitle>
                <CardDescription>Customize the look and feel</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-3">Color Theme</label>
                    <div className="flex gap-2">
                      <Button
                        variant={appearance.theme === 'light' ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setAppearance({ ...appearance, theme: 'light' })}
                        className="flex-1"
                      >
                        <Sun className="h-4 w-4 mr-2" />
                        Light
                      </Button>
                      <Button
                        variant={appearance.theme === 'dark' ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setAppearance({ ...appearance, theme: 'dark' })}
                        className="flex-1"
                      >
                        <Moon className="h-4 w-4 mr-2" />
                        Dark
                      </Button>
                      <Button
                        variant={appearance.theme === 'system' ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setAppearance({ ...appearance, theme: 'system' })}
                        className="flex-1"
                      >
                        <Monitor className="h-4 w-4 mr-2" />
                        System
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <SettingsIcon className="h-5 w-5" />
                  Layout
                </CardTitle>
                <CardDescription>Adjust layout preferences</CardDescription>
              </CardHeader>
              <CardContent>
                <SettingsRow title="Collapsed Sidebar" description="Start with sidebar collapsed">
                  <ToggleSwitch
                    checked={appearance.sidebarCollapsed}
                    onChange={(checked) => setAppearance({ ...appearance, sidebarCollapsed: checked })}
                    label="Toggle collapsed sidebar"
                  />
                </SettingsRow>
                <SettingsRow title="Compact Mode" description="Reduce spacing for more content">
                  <ToggleSwitch
                    checked={appearance.compactMode}
                    onChange={(checked) => setAppearance({ ...appearance, compactMode: checked })}
                    label="Toggle compact mode"
                  />
                </SettingsRow>
                <SettingsRow title="Show Avatars" description="Display user avatars in lists">
                  <ToggleSwitch
                    checked={appearance.showAvatars}
                    onChange={(checked) => setAppearance({ ...appearance, showAvatars: checked })}
                    label="Toggle show avatars"
                  />
                </SettingsRow>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* System Tab */}
        <TabsContent value="system" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Wifi className="h-5 w-5" />
                  Service Status
                </CardTitle>
                <CardDescription>Current status of system services</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between py-2 border-b border-gray-100">
                    <div className="flex items-center gap-2">
                      <HardDrive className="h-4 w-4 text-gray-400" />
                      <span className="text-sm">API Server</span>
                    </div>
                    <StatusIndicator status={systemStatus.api} label="Online" />
                  </div>
                  <div className="flex items-center justify-between py-2 border-b border-gray-100">
                    <div className="flex items-center gap-2">
                      <Database className="h-4 w-4 text-gray-400" />
                      <span className="text-sm">Database</span>
                    </div>
                    <StatusIndicator status={systemStatus.database} label="Online" />
                  </div>
                  <div className="flex items-center justify-between py-2 border-b border-gray-100">
                    <div className="flex items-center gap-2">
                      <RefreshCw className="h-4 w-4 text-gray-400" />
                      <span className="text-sm">Redis Cache</span>
                    </div>
                    <StatusIndicator status={systemStatus.redis} label="Online" />
                  </div>
                  <div className="flex items-center justify-between py-2 border-b border-gray-100">
                    <div className="flex items-center gap-2">
                      <Cloud className="h-4 w-4 text-gray-400" />
                      <span className="text-sm">NATS Messaging</span>
                    </div>
                    <StatusIndicator status={systemStatus.nats} label="Online" />
                  </div>
                  <div className="flex items-center justify-between py-2">
                    <div className="flex items-center gap-2">
                      <HardDrive className="h-4 w-4 text-gray-400" />
                      <span className="text-sm">MinIO Storage</span>
                    </div>
                    <StatusIndicator status={systemStatus.minio} label="Online" />
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <SettingsIcon className="h-5 w-5" />
                  System Information
                </CardTitle>
                <CardDescription>Application version and environment</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <label className="text-sm font-medium text-gray-700">Version</label>
                    <p className="text-sm text-gray-900">1.0.0</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-700">Environment</label>
                    <Badge className="bg-yellow-100 text-yellow-800 mt-1">Development</Badge>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-700">API Endpoint</label>
                    <p className="text-sm text-gray-900 font-mono bg-gray-50 p-2 rounded mt-1">
                      http://localhost:8080
                    </p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-700">Tenant ID</label>
                    <p className="text-sm text-gray-900 font-mono bg-gray-50 p-2 rounded mt-1">
                      demo-tenant
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Cloud className="h-5 w-5" />
                  SSOT Configuration
                </CardTitle>
                <CardDescription>Single Source of Truth sync settings</CardDescription>
              </CardHeader>
              <CardContent>
                <SettingsRow title="Auto Sync" description="Automatically sync with SSOT">
                  <ToggleSwitch
                    checked={ssotSettings.autoSync}
                    onChange={(checked) => setSsotSettings({ ...ssotSettings, autoSync: checked })}
                    label="Toggle auto sync"
                  />
                </SettingsRow>
                <SettingsRow title="Sync on Startup" description="Sync when dashboard loads">
                  <ToggleSwitch
                    checked={ssotSettings.syncOnStartup}
                    onChange={(checked) => setSsotSettings({ ...ssotSettings, syncOnStartup: checked })}
                    label="Toggle sync on startup"
                  />
                </SettingsRow>
                <div className="py-3 border-b border-gray-100">
                  <div className="flex items-center justify-between mb-2">
                    <div>
                      <p className="text-sm font-medium text-gray-900">Sync Interval</p>
                      <p className="text-xs text-gray-500">Minutes between auto syncs</p>
                    </div>
                    <Input
                      type="number"
                      min={5}
                      max={120}
                      value={ssotSettings.syncInterval}
                      onChange={(e) => setSsotSettings({ ...ssotSettings, syncInterval: parseInt(e.target.value) || 30 })}
                      className="w-20 text-center"
                    />
                  </div>
                </div>
                <div className="py-3">
                  <div className="flex items-center gap-2 text-sm text-gray-500">
                    <Clock className="h-4 w-4" />
                    Last sync: {new Date(lastSync).toLocaleString()}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  Data Management
                </CardTitle>
                <CardDescription>Manage application data</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between py-2">
                    <div>
                      <p className="text-sm font-medium text-gray-900">Export Data</p>
                      <p className="text-xs text-gray-500">Download all your data as CSV</p>
                    </div>
                    <Button variant="outline" size="sm" onClick={handleExportData}>
                      <Download className="h-4 w-4 mr-2" />
                      Export
                    </Button>
                  </div>
                  <div className="flex items-center justify-between py-2">
                    <div>
                      <p className="text-sm font-medium text-gray-900">Clear Cache</p>
                      <p className="text-xs text-gray-500">Clear local browser cache</p>
                    </div>
                    <Button variant="outline" size="sm" onClick={handleClearCache}>
                      <Trash2 className="h-4 w-4 mr-2" />
                      Clear
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Security Tab */}
        <TabsContent value="security" className="space-y-6">
          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="h-5 w-5" />
                  Session Information
                </CardTitle>
                <CardDescription>Current login session details</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <label className="text-sm font-medium text-gray-700">Session ID</label>
                    <p className="text-xs text-gray-900 font-mono bg-gray-50 p-2 rounded mt-1 truncate">
                      {Math.random().toString(36).substring(2, 15)}...
                    </p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-700">Login Time</label>
                    <p className="text-sm text-gray-900 mt-1">{new Date().toLocaleString()}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-700">IP Address</label>
                    <p className="text-sm text-gray-900 font-mono mt-1">127.0.0.1</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-700">User Agent</label>
                    <p className="text-xs text-gray-500 mt-1 truncate">
                      {typeof navigator !== 'undefined' ? navigator.userAgent.substring(0, 50) + '...' : 'N/A'}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Key className="h-5 w-5" />
                  Security Settings
                </CardTitle>
                <CardDescription>Manage account security</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between py-2 border-b border-gray-100">
                    <div>
                      <p className="text-sm font-medium text-gray-900">Two-Factor Authentication</p>
                      <p className="text-xs text-gray-500">Add an extra layer of security</p>
                    </div>
                    <Badge className="bg-gray-100 text-gray-600">Coming Soon</Badge>
                  </div>
                  <div className="flex items-center justify-between py-2 border-b border-gray-100">
                    <div>
                      <p className="text-sm font-medium text-gray-900">Active Sessions</p>
                      <p className="text-xs text-gray-500">View and manage active sessions</p>
                    </div>
                    <Button variant="outline" size="sm" disabled>
                      View All
                    </Button>
                  </div>
                  <div className="flex items-center justify-between py-2">
                    <div>
                      <p className="text-sm font-medium text-gray-900">Sign Out All Devices</p>
                      <p className="text-xs text-gray-500">End all active sessions</p>
                    </div>
                    <Button variant="outline" size="sm" className="text-red-600 hover:text-red-700">
                      <LogOut className="h-4 w-4 mr-2" />
                      Sign Out All
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
