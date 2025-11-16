import { apiClient } from './apiClient';
import { ENDPOINTS } from '../config/api.config';
import {
  SearchResponse,
  ContentDetailResponse,
  StatsResponse,
  SearchFilters,
} from '../types/content.types';

class ContentService {
  async searchContents(filters: SearchFilters): Promise<SearchResponse> {
    const params = new URLSearchParams();
    if (filters.q) params.append('q', filters.q);
    if (filters.type && filters.type !== 'all') params.append('type', filters.type);
    if (filters.sort) params.append('sort', filters.sort);
    if (filters.page) params.append('page', filters.page.toString());
    if (filters.page_size) params.append('page_size', filters.page_size.toString());
    const url = `${ENDPOINTS.CONTENTS.SEARCH}?${params.toString()}`;
    return apiClient.get<SearchResponse>(url);
  }

  async getContentById(id: string): Promise<ContentDetailResponse> {
    const url = ENDPOINTS.CONTENTS.GET_BY_ID(id);
    return apiClient.get<ContentDetailResponse>(url);
  }

  async getStats(): Promise<StatsResponse> {
    return apiClient.get<StatsResponse>(ENDPOINTS.CONTENTS.STATS);
  }
}

export const contentService = new ContentService();


