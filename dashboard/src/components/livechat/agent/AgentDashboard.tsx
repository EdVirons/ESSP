import { useState } from 'react';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useLivechatContext } from '@/contexts/LivechatContext';
import { AgentDashboardHeader } from './AgentDashboardHeader';
import { ChatSidebar } from './ChatSidebar';
import { ChatWindow } from './ChatWindow';
import { ChatMetricsBar } from './ChatMetricsBar';

export function AgentDashboard() {
  const { activeSessionId, setActiveSessionId } = useLivechatContext();
  const [mobileShowChat, setMobileShowChat] = useState(false);

  const handleChatAccepted = (sessionId: string) => {
    setActiveSessionId(sessionId);
    setMobileShowChat(true);
  };

  const handleBackToList = () => {
    setMobileShowChat(false);
  };

  return (
    <div className="flex flex-col h-full bg-gray-50">
      {/* Header */}
      <AgentDashboardHeader />

      {/* Mobile Metrics Bar */}
      <div className="lg:hidden px-4 py-2 bg-white border-b overflow-x-auto">
        <ChatMetricsBar />
      </div>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Desktop Layout: Side-by-side */}
        <div className="hidden md:flex flex-1">
          {/* Sidebar - 320px fixed width */}
          <div className="w-80 flex-shrink-0 border-r bg-white">
            <ChatSidebar onChatAccepted={handleChatAccepted} />
          </div>

          {/* Chat Window - Remaining space */}
          <div className="flex-1">
            <ChatWindow />
          </div>
        </div>

        {/* Mobile Layout: Either sidebar or chat */}
        <div className="flex md:hidden flex-1">
          {!mobileShowChat || !activeSessionId ? (
            // Show sidebar on mobile
            <div className="w-full bg-white">
              <ChatSidebar
                onChatAccepted={(sessionId) => {
                  handleChatAccepted(sessionId);
                }}
              />
            </div>
          ) : (
            // Show chat window on mobile
            <div className="w-full flex flex-col">
              {/* Back button */}
              <div className="px-4 py-2 bg-white border-b">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleBackToList}
                  className="text-gray-600"
                >
                  <ArrowLeft className="w-4 h-4 mr-1" />
                  Back to list
                </Button>
              </div>
              <div className="flex-1">
                <ChatWindow />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
