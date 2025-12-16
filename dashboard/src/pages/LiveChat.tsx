import { AgentDashboard } from '@/components/livechat/agent';
import { LivechatProvider } from '@/contexts/LivechatContext';

export default function LiveChat() {
  return (
    <LivechatProvider>
      <AgentDashboard />
    </LivechatProvider>
  );
}
