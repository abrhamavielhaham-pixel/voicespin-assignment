import { provideHttpClient } from '@angular/common/http';
import {
  HttpTestingController,
  provideHttpClientTesting,
} from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { Conversation } from './conversation';
import { ConversationService } from './conversation.service';

const resolved: Conversation = {
  id: 'conv-001',
  customerName: 'John Carter',
  customerEmail: 'john.carter@example.com',
  subject: 'Cannot reset my password',
  status: 'RESOLVED',
  priority: 'LOW',
  createdAt: '2026-09-01T09:15:00Z',
};

describe('ConversationService', () => {
  let service: ConversationService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });

    service = TestBed.inject(ConversationService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('lists conversations, sending only the filters that are set', () => {
    service.list({ search: '  john  ', status: 'OPEN', priority: '' }).subscribe();

    const req = httpMock.expectOne((request) => request.url === '/api/conversations');
    expect(req.request.method).toBe('GET');
    expect(req.request.params.get('search')).toBe('john');
    expect(req.request.params.get('status')).toBe('OPEN');
    expect(req.request.params.has('priority')).toBe(false);

    req.flush([]);
  });

  it('updates a conversation with PATCH and returns the updated resource', () => {
    let result: Conversation | undefined;
    service
      .update('conv-001', { status: 'RESOLVED', priority: 'LOW' })
      .subscribe((conversation) => (result = conversation));

    const req = httpMock.expectOne(
      (request) => request.url === '/api/conversations/conv-001',
    );
    expect(req.request.method).toBe('PATCH');
    expect(req.request.body).toEqual({ status: 'RESOLVED', priority: 'LOW' });

    req.flush(resolved);

    expect(result?.status).toBe('RESOLVED');
    expect(result?.priority).toBe('LOW');
  });
});
