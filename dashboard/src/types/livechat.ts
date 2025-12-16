import type { Message, MessageThread } from './messaging';

// Chat session status - includes AI active state
export type ChatSessionStatus = 'ai_active' | 'waiting' | 'active' | 'ended';

// Chat Session
export interface ChatSession {
  id: string;
  tenantId: string;
  schoolId: string;
  threadId: string;
  schoolContactId: string;
  schoolContactName: string;
  assignedAgentId?: string;
  assignedAgentName?: string;
  status: ChatSessionStatus;
  queuePosition?: number;
  startedAt: string;
  agentJoinedAt?: string;
  endedAt?: string;
  firstResponseSeconds?: number;
  totalMessages: number;
  rating?: number;
  feedback?: string;
  createdAt: string;
  updatedAt: string;
  thread?: MessageThread;

  // AI Support fields
  aiHandled: boolean;
  aiResolved?: boolean;
  aiTurns: number;
  escalationReason?: string;
  escalationSummary?: Record<string, unknown>;
  issueCategory?: string;
  issueSeverity?: string;
  collectedInfo?: Record<string, unknown>;
}

// Agent Availability
export interface AgentAvailability {
  tenantId: string;
  userId: string;
  isAvailable: boolean;
  maxConcurrentChats: number;
  currentChatCount: number;
  lastSeenAt: string;
  updatedAt: string;
  userName?: string;
}

// Request types
export interface StartSessionRequest {
  subject?: string;
  incidentId?: string;
}

export interface EndSessionRequest {
  rating?: number;
  feedback?: string;
}

export interface TransferChatRequest {
  targetAgentId: string;
  reason?: string;
}

export interface SetAvailabilityRequest {
  available: boolean;
  maxConcurrentChats?: number;
}

// Response types
export interface StartSessionResponse {
  session: ChatSession;
  thread: MessageThread;
  queuePosition?: number;
}

export interface QueuePositionResponse {
  position: number;
  estimatedWaitMinutes: number;
}

export interface AcceptChatResponse {
  session: ChatSession;
  thread: MessageThread;
}

export interface ChatQueueItem {
  sessionId: string;
  schoolId: string;
  schoolContactName: string;
  subject: string;
  waitingTimeSeconds: number;
  queuePosition: number;
  startedAt: string;
}

export interface ChatQueueResponse {
  items: ChatQueueItem[];
  totalWaiting: number;
  availableAgents: number;
}

export interface ActiveChatItem {
  session: ChatSession;
  thread: MessageThread;
  unreadCount: number;
  lastMessageAt?: string;
}

export interface ActiveChatsResponse {
  items: ActiveChatItem[];
  total: number;
}

export interface ChatMetrics {
  totalSessions: number;
  averageWaitTimeSeconds: number;
  averageResponseTimeSeconds: number;
  averageRating: number;
  sessionsWithRating: number;
  activeSessions: number;
  waitingSessions: number;
  endedSessions: number;
}

// WebSocket events for livechat
export interface ChatSessionUpdateEvent {
  type: 'chat_session_update';
  payload: {
    sessionId: string;
    status: ChatSessionStatus;
    agentId?: string;
    agentName?: string;
    queuePosition?: number;
  };
}

export interface PresenceUpdateEvent {
  type: 'presence_update';
  payload: {
    userId: string;
    status: 'online' | 'away' | 'offline';
  };
}

export interface NewChatWaitingEvent {
  type: 'new_chat_waiting';
  payload: {
    sessionId: string;
    schoolContactName: string;
    subject: string;
    startedAt: string;
    escalatedFromAI?: boolean;
    escalationReason?: string;
    aiSummary?: Record<string, unknown>;
  };
}

// AI Chat types

export interface AIChatMessageRequest {
  content: string;
  deviceSerial?: string;
}

export interface AIChatMessageResponse {
  message: Message;
  aiTyping: boolean;
  shouldEscalate: boolean;
  escalateReason?: string;
  sessionStatus: ChatSessionStatus;
}

export interface AIEscalationRequest {
  reason?: string;
}

export interface AIEscalationResponse {
  session: ChatSession;
  queuePosition?: number;
}

export interface AIConversationContext {
  sessionId: string;
  turnCount: number;
  category?: string;
  severity?: string;
  escalationReason?: string;
  summary?: Record<string, unknown>;
  collectedInfo?: Record<string, unknown>;
  conversationHistory: Message[];
  deviceContext?: Record<string, unknown>;
  schoolContext?: Record<string, unknown>;
}

// Typing indicator event
export interface TypingIndicatorEvent {
  type: 'typing_indicator';
  payload: {
    threadId: string;
    userId: string;
    userName: string;
    isTyping: boolean;
  };
}
