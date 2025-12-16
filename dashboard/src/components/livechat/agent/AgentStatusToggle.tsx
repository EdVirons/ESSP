import { useState } from 'react';
import { Circle, Settings, ChevronDown } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Popover } from '@/components/ui/popover';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { useAgentAvailability, useSetAvailability } from '@/hooks/useLivechat';
import { useLivechatContext } from '@/contexts/LivechatContext';

export function AgentStatusToggle() {
  const { data: availability, isLoading } = useAgentAvailability();
  const setAvailability = useSetAvailability();
  const { setAgentAvailable } = useLivechatContext();
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [maxChats, setMaxChats] = useState(availability?.maxConcurrentChats || 5);

  const isOnline = availability?.isAvailable ?? false;
  const currentCount = availability?.currentChatCount ?? 0;

  const handleToggle = () => {
    const newStatus = !isOnline;
    setAvailability.mutate(
      { available: newStatus, maxConcurrentChats: maxChats },
      {
        onSuccess: () => {
          setAgentAvailable(newStatus);
        },
      }
    );
  };

  const handleMaxChatsChange = (value: string) => {
    const newMax = Math.max(1, Math.min(10, parseInt(value) || 1));
    setMaxChats(newMax);
    if (isOnline) {
      setAvailability.mutate({ available: true, maxConcurrentChats: newMax });
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-100 animate-pulse">
        <div className="w-3 h-3 rounded-full bg-gray-300" />
        <div className="w-16 h-4 bg-gray-300 rounded" />
      </div>
    );
  }

  return (
    <div className="flex items-center gap-3">
      {/* Status Toggle Button */}
      <button
        onClick={handleToggle}
        disabled={setAvailability.isPending}
        className={`flex items-center gap-2 px-3 py-2 rounded-lg transition-colors ${
          isOnline
            ? 'bg-green-50 border border-green-200 hover:bg-green-100'
            : 'bg-gray-50 border border-gray-200 hover:bg-gray-100'
        }`}
      >
        <Circle
          className={`w-3 h-3 ${
            isOnline ? 'fill-green-500 text-green-500' : 'fill-gray-400 text-gray-400'
          }`}
        />
        <span
          className={`text-sm font-medium ${
            isOnline ? 'text-green-700' : 'text-gray-600'
          }`}
        >
          {isOnline ? 'Online' : 'Offline'}
        </span>
        <ChevronDown className="w-4 h-4 text-gray-400" />
      </button>

      {isOnline && (
        <span className="text-sm text-gray-500">
          {currentCount}/{maxChats} chats
        </span>
      )}

      {/* Settings Popover */}
      <Popover
        open={settingsOpen}
        onClose={() => setSettingsOpen(false)}
        trigger={
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => setSettingsOpen(!settingsOpen)}
          >
            <Settings className="h-4 w-4 text-gray-500" />
          </Button>
        }
        align="end"
      >
        <div className="p-4 space-y-4">
          <div>
            <h4 className="font-medium text-sm mb-2">Chat Settings</h4>
          </div>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="max-chats" className="text-sm">
                Max concurrent chats
              </Label>
            </div>
            <Input
              id="max-chats"
              type="number"
              min={1}
              max={10}
              value={maxChats}
              onChange={(e) => handleMaxChatsChange(e.target.value)}
              className="w-full"
            />
            <p className="text-xs text-gray-500">
              New chats won't be assigned when you reach this limit (1-10)
            </p>
          </div>
        </div>
      </Popover>
    </div>
  );
}
