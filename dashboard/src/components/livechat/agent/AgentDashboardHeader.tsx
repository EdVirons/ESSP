import { Headphones } from 'lucide-react';
import { AgentStatusToggle } from './AgentStatusToggle';
import { ChatMetricsBar } from './ChatMetricsBar';

export function AgentDashboardHeader() {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 px-6 py-4 bg-white border-b">
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-cyan-100">
          <Headphones className="h-5 w-5 text-cyan-600" />
        </div>
        <div>
          <h1 className="text-lg font-semibold text-gray-900">Live Chat</h1>
          <p className="text-sm text-gray-500">Manage customer conversations</p>
        </div>
      </div>

      <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
        <AgentStatusToggle />
        <div className="hidden lg:block">
          <ChatMetricsBar />
        </div>
      </div>
    </div>
  );
}
