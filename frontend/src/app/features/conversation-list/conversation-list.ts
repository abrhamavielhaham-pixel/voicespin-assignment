import { ChangeDetectionStrategy, Component, input, output, signal } from '@angular/core';
import { DatePipe } from '@angular/common';
import {
  Conversation,
  ConversationFilters,
  ConversationPriority,
  ConversationStatus,
  PRIORITIES,
  STATUSES,
  STATUS_LABELS,
} from '../../core/conversation';

@Component({
  selector: 'app-conversation-list',
  imports: [DatePipe],
  templateUrl: './conversation-list.html',
  styleUrl: './conversation-list.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ConversationList {
  readonly conversations = input.required<Conversation[]>();
  readonly selectedId = input<string | null>(null);
  readonly loading = input(false);
  readonly error = input<string | null>(null);

  readonly selectConversation = output<string>();
  readonly filtersChange = output<ConversationFilters>();

  protected readonly statuses = STATUSES;
  protected readonly priorities = PRIORITIES;
  protected readonly statusLabels = STATUS_LABELS;

  protected readonly search = signal('');
  protected readonly status = signal<ConversationStatus | ''>('');
  protected readonly priority = signal<ConversationPriority | ''>('');

  protected onSearchInput(value: string): void {
    this.search.set(value);
    this.emitFilters();
  }

  protected onStatusChange(value: string): void {
    this.status.set(value as ConversationStatus | '');
    this.emitFilters();
  }

  protected onPriorityChange(value: string): void {
    this.priority.set(value as ConversationPriority | '');
    this.emitFilters();
  }

  private emitFilters(): void {
    this.filtersChange.emit({
      search: this.search(),
      status: this.status(),
      priority: this.priority(),
    });
  }
}
