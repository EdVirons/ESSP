export type NotificationType = 'incident' | 'work_order' | 'project' | 'device' | 'service_shop';
export type NotificationAction = 'create' | 'update' | 'delete';

export interface Notification {
  id: string;
  type: NotificationType;
  action: NotificationAction;
  actor: string;
  target: string;
  summary: string;
  timestamp: string;
  metadata: Record<string, unknown>;
  read: boolean;
}

export interface NotificationsResponse {
  items: Notification[];
  unreadCount: number;
  total: number;
}

export interface UnreadCountResponse {
  count: number;
}

export interface MarkReadRequest {
  ids: string; // comma-separated IDs or "all"
}

// WebSocket message types
export type WSMessageType =
  | 'notification'
  | 'entity_update'
  | 'ping'
  | 'pong'
  // Chat/messaging types
  | 'chat_message'
  | 'chat_typing'
  | 'chat_read'
  | 'chat_session_update'
  | 'presence_update'
  | 'new_chat_waiting'
  | 'agent_assigned'
  | 'typing_indicator';

export interface WSMessage<T = unknown> {
  type: WSMessageType;
  payload: T;
  timestamp?: string;
}
