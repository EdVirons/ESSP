// Thread types
export type ThreadType = 'general' | 'incident' | 'livechat';
export type ThreadStatus = 'open' | 'closed' | 'archived';
export type ContentType = 'text' | 'system' | 'attachment';

// Message Thread
export interface MessageThread {
  id: string;
  tenantId: string;
  schoolId: string;
  subject: string;
  threadType: ThreadType;
  status: ThreadStatus;
  incidentId?: string;
  createdBy: string;
  createdByRole: string;
  createdByName: string;
  messageCount: number;
  unreadCountSchool: number;
  unreadCountSupport: number;
  lastMessageAt?: string;
  createdAt: string;
  updatedAt: string;
  closedAt?: string;
  schoolName?: string;
  lastMessage?: Message;
}

// Message
export interface Message {
  id: string;
  tenantId: string;
  threadId: string;
  senderId: string;
  senderName: string;
  senderRole: string;
  content: string;
  contentType: ContentType;
  metadata?: Record<string, unknown>;
  editedAt?: string;
  deletedAt?: string;
  createdAt: string;
  attachments?: MessageAttachment[];
}

// Message Attachment
export interface MessageAttachment {
  id: string;
  tenantId: string;
  messageId: string;
  fileName: string;
  contentType: string;
  sizeBytes: number;
  objectKey: string;
  thumbnailKey?: string;
  downloadUrl?: string;
  thumbnailUrl?: string;
  createdAt: string;
}

// Thread Participant
export interface ThreadParticipant {
  threadId: string;
  userId: string;
  userName: string;
  userRole: string;
  joinedAt: string;
  leftAt?: string;
}

// Read Receipt
export interface ReadReceipt {
  threadId: string;
  userId: string;
  lastReadMessageId: string;
  lastReadAt: string;
}

// Unread Counts
export interface UnreadCounts {
  threads: number;
  messages: number;
  total: number;
}

// Request types
export interface CreateThreadRequest {
  subject: string;
  incidentId?: string;
  initialMessage: string;
  attachments?: string[];
  schoolId?: string;
}

export interface CreateMessageRequest {
  content: string;
  attachments?: string[];
}

export interface UpdateThreadStatusRequest {
  status: ThreadStatus;
}

export interface MarkReadRequest {
  lastMessageId: string;
}

export interface UploadURLRequest {
  fileName: string;
  contentType: string;
  sizeBytes: number;
}

// Response types
export interface ThreadsListResponse {
  items: MessageThread[];
  nextCursor?: string;
}

export interface MessagesListResponse {
  items: Message[];
  nextCursor?: string;
}

export interface ThreadDetailResponse {
  thread: MessageThread;
  messages: Message[];
  participants?: ThreadParticipant[];
}

export interface CreateThreadResponse {
  thread: MessageThread;
  message: Message;
}

export interface CreateMessageResponse {
  message: Message;
}

export interface UploadURLResponse {
  uploadUrl: string;
  attachmentRef: string;
}

export interface MessageSearchResult {
  message: Message;
  thread: MessageThread;
  highlights?: string[];
}

export interface SearchResponse {
  items: MessageSearchResult[];
  total: number;
}

// Analytics
export interface MessagingAnalytics {
  threadsCreated: number;
  messagesSent: number;
  avgResponseTimeMinutes: number;
  chatSessions: number;
  avgChatRating: number;
  activeThreads: number;
  closedThreads: number;
  avgMessagesPerThread: number;
}

// WebSocket message types for real-time updates
export interface ChatMessageEvent {
  type: 'chat_message';
  payload: {
    threadId: string;
    message: Message;
  };
}

export interface ChatTypingEvent {
  type: 'chat_typing';
  payload: {
    threadId: string;
    userId: string;
    userName: string;
    isTyping: boolean;
  };
}

export interface ChatReadEvent {
  type: 'chat_read';
  payload: {
    threadId: string;
    userId: string;
    lastMessageId: string;
    readAt: string;
  };
}

// List params
export interface ThreadsListParams {
  status?: string;
  incidentId?: string;
  q?: string;
  cursor?: string;
  limit?: number;
}
