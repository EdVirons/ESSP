import { MessagesPage } from '@/components/messaging/MessagesPage';

export function Messages() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Messages</h1>
        <p className="text-gray-500">
          Communicate with school contacts and support team
        </p>
      </div>

      <MessagesPage />
    </div>
  );
}
