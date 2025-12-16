import { useState } from 'react';
import { ArrowRightLeft, User } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { useTransferChat } from '@/hooks/useLivechat';
import { useLivechatContext } from '@/contexts/LivechatContext';

interface ChatTransferModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  sessionId: string;
}

export function ChatTransferModal({ open, onOpenChange, sessionId }: ChatTransferModalProps) {
  const [targetAgentId, setTargetAgentId] = useState('');
  const [reason, setReason] = useState('');
  const transferChat = useTransferChat();
  const { setActiveSessionId } = useLivechatContext();

  const handleTransfer = () => {
    if (!targetAgentId.trim()) return;

    transferChat.mutate(
      {
        sessionId,
        data: {
          targetAgentId: targetAgentId.trim(),
          reason: reason.trim() || undefined,
        },
      },
      {
        onSuccess: () => {
          onOpenChange(false);
          setTargetAgentId('');
          setReason('');
          setActiveSessionId(null);
        },
      }
    );
  };

  const handleClose = () => {
    onOpenChange(false);
    setTargetAgentId('');
    setReason('');
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ArrowRightLeft className="w-5 h-5 text-cyan-600" />
            Transfer Chat
          </DialogTitle>
          <DialogDescription>
            Transfer this chat to another support agent. The conversation history will be preserved.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="agent-id">Target Agent ID</Label>
            <div className="relative">
              <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <Input
                id="agent-id"
                value={targetAgentId}
                onChange={(e) => setTargetAgentId(e.target.value)}
                placeholder="Enter agent ID or username"
                className="pl-10"
              />
            </div>
            <p className="text-xs text-gray-500">
              Enter the ID of the agent you want to transfer to
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="reason">Reason (optional)</Label>
            <Textarea
              id="reason"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Why are you transferring this chat?"
              rows={3}
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            onClick={handleTransfer}
            disabled={!targetAgentId.trim() || transferChat.isPending}
            className="bg-cyan-600 hover:bg-cyan-700"
          >
            {transferChat.isPending ? 'Transferring...' : 'Transfer'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
