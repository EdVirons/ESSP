import { useState } from 'react';
import { Loader2, Send } from 'lucide-react';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useCreateThread } from '@/hooks/useMessages';

interface NewThreadModalProps {
  isOpen: boolean;
  onClose: () => void;
  onThreadCreated?: (threadId: string) => void;
  incidentId?: string;
  schoolId?: string;
}

export function NewThreadModal({
  isOpen,
  onClose,
  onThreadCreated,
  incidentId,
}: NewThreadModalProps) {
  const [subject, setSubject] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState<string | null>(null);

  const createThread = useCreateThread();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!subject.trim() || !message.trim()) {
      setError('Please fill in all fields');
      return;
    }

    try {
      const result = await createThread.mutateAsync({
        subject: subject.trim(),
        initialMessage: message.trim(),
        incidentId,
      });

      // Reset form
      setSubject('');
      setMessage('');
      onClose();
      onThreadCreated?.(result.thread.id);
    } catch {
      setError('Failed to start conversation');
    }
  };

  const handleClose = () => {
    setSubject('');
    setMessage('');
    setError(null);
    onClose();
  };

  return (
    <Modal open={isOpen} onClose={handleClose} className="max-w-[500px]">
      <ModalHeader onClose={handleClose}>New Conversation</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody className="space-y-4">
          {error && (
            <div className="p-3 text-sm text-red-600 bg-red-50 rounded-lg">
              {error}
            </div>
          )}

          <div className="space-y-2">
            <label htmlFor="subject" className="block text-sm font-medium text-gray-700">
              Subject
            </label>
            <Input
              id="subject"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              placeholder="What is this conversation about?"
              required
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="message" className="block text-sm font-medium text-gray-700">
              Message
            </label>
            <Textarea
              id="message"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Write your message..."
              rows={5}
              required
            />
          </div>

          {incidentId && (
            <div className="text-sm text-gray-500 bg-gray-50 p-3 rounded-lg">
              This conversation will be linked to incident{' '}
              <span className="font-mono font-medium">{incidentId}</span>
            </div>
          )}
        </ModalBody>

        <ModalFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={createThread.isPending}
            className="bg-cyan-600 hover:bg-cyan-700"
          >
            {createThread.isPending ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Sending...
              </>
            ) : (
              <>
                <Send className="h-4 w-4 mr-2" />
                Start Conversation
              </>
            )}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
