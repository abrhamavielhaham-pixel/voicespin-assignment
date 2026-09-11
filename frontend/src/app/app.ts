import {
  ChangeDetectionStrategy,
  Component,
  computed,
  inject,
  signal,
} from '@angular/core';
import {
  Conversation,
  ConversationFilters,
  ConversationUpdate,
} from './core/conversation';
import { ConversationService } from './core/conversation.service';
import { ConversationDetail } from './features/conversation-detail/conversation-detail';
import { ConversationList } from './features/conversation-list/conversation-list';

@Component({
  selector: 'app-root',
  imports: [ConversationList, ConversationDetail],
  templateUrl: './app.html',
  styleUrl: './app.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class App {
  private readonly service = inject(ConversationService);

  protected readonly filters = signal<ConversationFilters>({
    search: '',
    status: '',
    priority: '',
  });
  protected readonly conversations = signal<Conversation[]>([]);
  protected readonly selectedId = signal<string | null>(null);
  protected readonly loading = signal(false);
  protected readonly loadError = signal<string | null>(null);
  protected readonly saving = signal(false);
  protected readonly saveError = signal<string | null>(null);

  protected readonly selected = computed(
    () => this.conversations().find((c) => c.id === this.selectedId()) ?? null,
  );

  constructor() {
    this.load();
  }

  protected onFiltersChange(filters: ConversationFilters): void {
    this.filters.set(filters);
    this.load();
  }

  protected onSelect(id: string): void {
    this.selectedId.set(id);
    this.saveError.set(null);
  }

  protected onSave(changes: ConversationUpdate): void {
    const id = this.selectedId();
    if (!id) {
      return;
    }

    this.saving.set(true);
    this.saveError.set(null);
    this.service.update(id, changes).subscribe({
      next: (updated) => {
        this.conversations.update((list) =>
          list.map((c) => (c.id === updated.id ? updated : c)),
        );
        this.saving.set(false);
      },
      error: () => {
        this.saving.set(false);
        this.saveError.set('Could not save changes. Please try again.');
      },
    });
  }

  private load(): void {
    this.loading.set(true);
    this.loadError.set(null);
    this.service.list(this.filters()).subscribe({
      next: (conversations) => {
        this.conversations.set(conversations);
        this.loading.set(false);
        if (
          this.selectedId() &&
          !conversations.some((c) => c.id === this.selectedId())
        ) {
          this.selectedId.set(null);
        }
      },
      error: () => {
        this.loading.set(false);
        this.loadError.set('Could not load conversations. Is the backend running?');
      },
    });
  }
}
