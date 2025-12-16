import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ChatQueue } from './ChatQueue';
import { ActiveChats } from './ActiveChats';
import { useChatQueue, useActiveChats } from '@/hooks/useLivechat';

interface ChatSidebarProps {
  onChatAccepted?: (sessionId: string) => void;
}

export function ChatSidebar({ onChatAccepted }: ChatSidebarProps) {
  const [activeTab, setActiveTab] = useState<'queue' | 'active'>('queue');
  const { data: queue } = useChatQueue();
  const { data: activeChats } = useActiveChats();

  const queueCount = queue?.totalWaiting ?? 0;
  const activeCount = activeChats?.total ?? 0;

  return (
    <div className="flex flex-col h-full bg-white border-r">
      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'queue' | 'active')} className="flex flex-col h-full">
        <div className="border-b px-2 pt-2">
          <TabsList className="w-full grid grid-cols-2">
            <TabsTrigger value="queue" className="relative">
              Queue
              {queueCount > 0 && (
                <span className="ml-1.5 px-1.5 py-0.5 text-xs bg-amber-500 text-white rounded-full">
                  {queueCount}
                </span>
              )}
            </TabsTrigger>
            <TabsTrigger value="active" className="relative">
              Active
              {activeCount > 0 && (
                <span className="ml-1.5 px-1.5 py-0.5 text-xs bg-cyan-500 text-white rounded-full">
                  {activeCount}
                </span>
              )}
            </TabsTrigger>
          </TabsList>
        </div>
        <TabsContent value="queue" className="flex-1 m-0 overflow-hidden">
          <ChatQueue onChatAccepted={onChatAccepted} />
        </TabsContent>
        <TabsContent value="active" className="flex-1 m-0 overflow-hidden">
          <ActiveChats />
        </TabsContent>
      </Tabs>
    </div>
  );
}
