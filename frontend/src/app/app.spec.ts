import { provideHttpClient } from '@angular/common/http';
import {
  HttpTestingController,
  provideHttpClientTesting,
} from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { App } from './app';
import { Conversation } from './core/conversation';

const john: Conversation = {
  id: 'conv-001',
  customerName: 'John Carter',
  customerEmail: 'john.carter@example.com',
  subject: 'Cannot reset my password',
  status: 'OPEN',
  priority: 'HIGH',
  createdAt: '2026-09-01T09:15:00Z',
};

const maria: Conversation = {
  id: 'conv-002',
  customerName: 'Maria Lopez',
  customerEmail: 'maria.lopez@example.com',
  subject: 'Invoice charged twice',
  status: 'OPEN',
  priority: 'LOW',
  createdAt: '2026-09-02T11:40:00Z',
};

describe('App', () => {
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [App],
      providers: [provideHttpClient(), provideHttpClientTesting()],
    }).compileComponents();

    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  function root(fixture: ComponentFixture<App>): HTMLElement {
    return fixture.nativeElement as HTMLElement;
  }

  async function render(conversations: Conversation[]): Promise<ComponentFixture<App>> {
    const fixture = TestBed.createComponent(App);

    httpMock.expectOne('/api/conversations').flush(conversations);

    await fixture.whenStable();
    fixture.detectChanges();
    return fixture;
  }

  function rows(fixture: ComponentFixture<App>): HTMLElement[] {
    return Array.from(root(fixture).querySelectorAll<HTMLElement>('.conversation-item'));
  }

  async function selectFirst(fixture: ComponentFixture<App>): Promise<void> {
    rows(fixture)[0].click();
    await fixture.whenStable();
    fixture.detectChanges();
  }

  it('renders the conversations returned by the API', async () => {
    const fixture = await render([john, maria]);

    expect(rows(fixture).length).toBe(2);
    expect(rows(fixture)[0].textContent).toContain('John Carter');
    expect(rows(fixture)[0].textContent).toContain('Cannot reset my password');
    expect(rows(fixture)[0].textContent).toContain('Open');
  });

  it('shows an empty state when nothing matches', async () => {
    const fixture = await render([]);

    expect(rows(fixture).length).toBe(0);
    expect(root(fixture).querySelector('.empty-state')?.textContent).toContain(
      'No conversations match your filters.',
    );
  });

  it('shows an error state when the request fails', async () => {
    const fixture = TestBed.createComponent(App);

    httpMock
      .expectOne('/api/conversations')
      .flush({ error: 'boom' }, { status: 500, statusText: 'Server Error' });

    await fixture.whenStable();
    fixture.detectChanges();

    expect(root(fixture).querySelector('.error-state')?.textContent).toContain(
      'Could not load conversations',
    );
  });

  it('reloads with a search query when the user types', async () => {
    const fixture = await render([john, maria]);

    const searchInput = root(fixture).querySelector<HTMLInputElement>('.search-input')!;
    searchInput.value = 'john';
    searchInput.dispatchEvent(new Event('input'));

    httpMock.expectOne((req) => req.params.get('search') === 'john').flush([john]);

    await fixture.whenStable();
    fixture.detectChanges();

    expect(rows(fixture).length).toBe(1);
    expect(rows(fixture)[0].textContent).toContain('John Carter');
  });

  it('shows details for the selected conversation', async () => {
    const fixture = await render([john, maria]);

    await selectFirst(fixture);

    const detail = root(fixture).querySelector('app-conversation-detail');
    expect(detail?.textContent).toContain('john.carter@example.com');
    expect(
      root(fixture).querySelector<HTMLSelectElement>('#detail-status')?.value,
    ).toBe('OPEN');
    expect(
      root(fixture).querySelector<HTMLSelectElement>('#detail-priority')?.value,
    ).toBe('HIGH');
  });

  it('saves status and priority changes through the API', async () => {
    const fixture = await render([john, maria]);

    await selectFirst(fixture);

    const statusSelect = root(fixture).querySelector<HTMLSelectElement>('#detail-status')!;
    statusSelect.value = 'RESOLVED';
    statusSelect.dispatchEvent(new Event('change'));

    const prioritySelect =
      root(fixture).querySelector<HTMLSelectElement>('#detail-priority')!;
    prioritySelect.value = 'LOW';
    prioritySelect.dispatchEvent(new Event('change'));

    fixture.detectChanges();
    root(fixture).querySelector<HTMLButtonElement>('.save-button')!.click();

    const req = httpMock.expectOne('/api/conversations/conv-001');
    expect(req.request.method).toBe('PATCH');
    expect(req.request.body).toEqual({ status: 'RESOLVED', priority: 'LOW' });
    req.flush({ ...john, status: 'RESOLVED', priority: 'LOW' });

    await fixture.whenStable();
    fixture.detectChanges();

    expect(rows(fixture)[0].querySelector('.status-resolved')?.textContent).toContain(
      'Resolved',
    );
    expect(rows(fixture)[0].querySelector('.priority-low')?.textContent).toContain('LOW');
  });
});
