import { Users, MessageSquare, Clock } from 'lucide-react';
import { useChatQueue, useActiveChats, useChatMetrics } from '@/hooks/useLivechat';

export function ChatMetricsBar() {
  const { data: queue } = useChatQueue();
  const { data: activeChats } = useActiveChats();
  const { data: metrics } = useChatMetrics();

  const queueCount = queue?.totalWaiting ?? 0;
  const activeCount = activeChats?.total ?? 0;
  const avgWaitSeconds = metrics?.averageWaitTimeSeconds ?? 0;

  const formatWaitTime = (seconds: number): string => {
    if (seconds < 60) return `${Math.round(seconds)}s`;
    const minutes = Math.round(seconds / 60);
    return `${minutes}m`;
  };

  return (
    <div className="flex items-center gap-4">
      <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-amber-50 border border-amber-200">
        <Users className="w-4 h-4 text-amber-600" />
        <span className="text-sm font-medium text-amber-700">
          Queue: {queueCount}
        </span>
      </div>

      <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-cyan-50 border border-cyan-200">
        <MessageSquare className="w-4 h-4 text-cyan-600" />
        <span className="text-sm font-medium text-cyan-700">
          Active: {activeCount}
        </span>
      </div>

      <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-gray-50 border border-gray-200">
        <Clock className="w-4 h-4 text-gray-600" />
        <span className="text-sm font-medium text-gray-700">
          Avg Wait: {formatWaitTime(avgWaitSeconds)}
        </span>
      </div>
    </div>
  );
}
