import { Settings, Bell, Globe, Moon, Sun, Monitor, Clock, Check, X } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import type { UserPreferences } from '@/lib/api';

interface UserPreferencesCardProps {
  preferences?: UserPreferences;
  onEdit?: () => void;
}

// Theme icon mapping
const themeIcons: Record<string, React.ReactNode> = {
  light: <Sun className="h-4 w-4 text-amber-500" />,
  dark: <Moon className="h-4 w-4 text-indigo-500" />,
  system: <Monitor className="h-4 w-4 text-gray-500" />,
};

export function UserPreferencesCard({ preferences, onEdit }: UserPreferencesCardProps) {
  const theme = preferences?.theme || 'light';
  const language = preferences?.language || 'en';
  const timezone = preferences?.timezone || 'America/New_York';
  const notifications = preferences?.notifications;

  return (
    <Card className="border-0 shadow-md">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-purple-500 to-indigo-600">
                <Settings className="h-4 w-4 text-white" />
              </div>
              Preferences
            </CardTitle>
            <CardDescription>
              Your dashboard preferences
            </CardDescription>
          </div>
          {onEdit && (
            <Button variant="outline" size="sm" onClick={onEdit}>
              Edit
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {/* Appearance */}
        <div className="flex items-center justify-between rounded-xl bg-gradient-to-r from-amber-50 to-orange-50 border border-amber-100 p-3">
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white shadow-sm">
              {themeIcons[theme]}
            </div>
            <span className="text-sm font-medium text-gray-700">Theme</span>
          </div>
          <span className="text-sm font-semibold text-gray-900 capitalize bg-white px-2 py-0.5 rounded-md shadow-sm">{theme}</span>
        </div>

        {/* Language */}
        <div className="flex items-center justify-between rounded-xl bg-gradient-to-r from-blue-50 to-cyan-50 border border-blue-100 p-3">
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white shadow-sm">
              <Globe className="h-4 w-4 text-blue-500" />
            </div>
            <span className="text-sm font-medium text-gray-700">Language</span>
          </div>
          <span className="text-sm font-semibold text-gray-900 bg-white px-2 py-0.5 rounded-md shadow-sm">
            {language === 'en' ? 'English' : language === 'es' ? 'Spanish' : language}
          </span>
        </div>

        {/* Timezone */}
        <div className="flex items-center justify-between rounded-xl bg-gradient-to-r from-emerald-50 to-teal-50 border border-emerald-100 p-3">
          <div className="flex items-center gap-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white shadow-sm">
              <Clock className="h-4 w-4 text-emerald-500" />
            </div>
            <span className="text-sm font-medium text-gray-700">Timezone</span>
          </div>
          <span className="text-sm font-semibold text-gray-900 bg-white px-2 py-0.5 rounded-md shadow-sm text-xs">{timezone}</span>
        </div>

        {/* Notifications */}
        {notifications && (
          <div className="rounded-xl bg-gradient-to-r from-purple-50 to-pink-50 border border-purple-100 p-3 space-y-3">
            <div className="flex items-center gap-3">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white shadow-sm">
                <Bell className="h-4 w-4 text-purple-500" />
              </div>
              <span className="text-sm font-medium text-gray-700">Notifications</span>
            </div>
            <div className="grid grid-cols-2 gap-2">
              <div className="flex items-center justify-between text-xs bg-white rounded-lg p-2 shadow-sm">
                <span className="text-gray-600">Email</span>
                {notifications.emailEnabled ? (
                  <span className="flex items-center gap-1 text-emerald-600">
                    <Check className="h-3 w-3" /> On
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-gray-400">
                    <X className="h-3 w-3" /> Off
                  </span>
                )}
              </div>
              <div className="flex items-center justify-between text-xs bg-white rounded-lg p-2 shadow-sm">
                <span className="text-gray-600">Browser</span>
                {notifications.browserEnabled ? (
                  <span className="flex items-center gap-1 text-emerald-600">
                    <Check className="h-3 w-3" /> On
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-gray-400">
                    <X className="h-3 w-3" /> Off
                  </span>
                )}
              </div>
              <div className="flex items-center justify-between text-xs bg-white rounded-lg p-2 shadow-sm">
                <span className="text-gray-600">Incidents</span>
                {notifications.incidentAlerts ? (
                  <span className="flex items-center gap-1 text-emerald-600">
                    <Check className="h-3 w-3" /> On
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-gray-400">
                    <X className="h-3 w-3" /> Off
                  </span>
                )}
              </div>
              <div className="flex items-center justify-between text-xs bg-white rounded-lg p-2 shadow-sm">
                <span className="text-gray-600">Work Orders</span>
                {notifications.workOrderAlerts ? (
                  <span className="flex items-center gap-1 text-emerald-600">
                    <Check className="h-3 w-3" /> On
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-gray-400">
                    <X className="h-3 w-3" /> Off
                  </span>
                )}
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
