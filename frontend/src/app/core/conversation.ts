export type ConversationStatus = 'OPEN' | 'IN_PROGRESS' | 'RESOLVED';
export type ConversationPriority = 'LOW' | 'MEDIUM' | 'HIGH';

export interface Conversation {
  id: string;
  customerName: string;
  customerEmail: string;
  subject: string;
  status: ConversationStatus;
  priority: ConversationPriority;
  createdAt: string;
}

export interface ConversationFilters {
  search: string;
  status: ConversationStatus | '';
  priority: ConversationPriority | '';
}

export interface ConversationUpdate {
  status?: ConversationStatus;
  priority?: ConversationPriority;
}

export const STATUSES: readonly ConversationStatus[] = ['OPEN', 'IN_PROGRESS', 'RESOLVED'];
export const PRIORITIES: readonly ConversationPriority[] = ['LOW', 'MEDIUM', 'HIGH'];

export const STATUS_LABELS: Record<ConversationStatus, string> = {
  OPEN: 'Open',
  IN_PROGRESS: 'In progress',
  RESOLVED: 'Resolved',
};
