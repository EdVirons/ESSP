import { AlertTriangle } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { useEndSession } from '@/hooks/useLivechat';
import { useLivechatContext } from '@/contexts/LivechatContext';

interface EndChatModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  sessionId: string;
  contactName: string;
}

export function EndChatModal({ open, onOpenChange, sessionId, contactName }: EndChatModalProps) {
  const endSession = useEndSession();
  const { setActiveSessionId } = useLivechatContext();

  const handleEnd = () => {
    endSession.mutate(
      { sessionId },
      {
        onSuccess: () => {
          onOpenChange(false);
          setActiveSessionId(null);
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <AlertTriangle className="w-5 h-5 text-amber-500" />
            End Chat Session
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to end the chat with <strong>{contactName}</strong>?
            The school contact will be notified that the session has ended.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleEnd}
            disabled={endSession.isPending}
            className="bg-red-600 hover:bg-red-700"
          >
            {endSession.isPending ? 'Ending...' : 'End Chat'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
