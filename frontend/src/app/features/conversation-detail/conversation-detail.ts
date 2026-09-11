import {
  ChangeDetectionStrategy,
  Component,
  computed,
  input,
  linkedSignal,
  output,
} from '@angular/core';
import { DatePipe } from '@angular/common';
import {
  Conversation,
  ConversationUpdate,
  PRIORITIES,
  STATUSES,
  STATUS_LABELS,
} from '../../core/conversation';

@Component({
  selector: 'app-conversation-detail',
  imports: [DatePipe],
  templateUrl: './conversation-detail.html',
  styleUrl: './conversation-detail.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ConversationDetail {
  readonly conversation = input.required<Conversation>();
  readonly saving = input(false);
  readonly saveError = input<string | null>(null);

  readonly save = output<ConversationUpdate>();

  protected readonly statuses = STATUSES;
  protected readonly priorities = PRIORITIES;
  protected readonly statusLabels = STATUS_LABELS;

  protected readonly status = linkedSignal(() => this.conversation().status);
  protected readonly priority = linkedSignal(() => this.conversation().priority);

  protected readonly dirty = computed(
    () =>
      this.status() !== this.conversation().status ||
      this.priority() !== this.conversation().priority,
  );

  protected onSave(): void {
    if (!this.dirty()) {
      return;
    }
    this.save.emit({ status: this.status(), priority: this.priority() });
  }
}
