import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { Conversation, ConversationFilters, ConversationUpdate } from './conversation';

@Injectable({ providedIn: 'root' })
export class ConversationService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/conversations';

  list(filters: ConversationFilters): Observable<Conversation[]> {
    let params = new HttpParams();
    if (filters.status) {
      params = params.set('status', filters.status);
    }
    if (filters.priority) {
      params = params.set('priority', filters.priority);
    }
    const search = filters.search.trim();
    if (search) {
      params = params.set('search', search);
    }
    return this.http.get<Conversation[]>(this.baseUrl, { params });
  }

  update(id: string, changes: ConversationUpdate): Observable<Conversation> {
    return this.http.patch<Conversation>(`${this.baseUrl}/${id}`, changes);
  }
}
