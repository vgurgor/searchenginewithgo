export type ContentType = 'video' | 'text';

export interface Content {
  id: string;
  title: string;
  content_type: ContentType;
  description: string;
  url?: string;
  thumbnail_url?: string;
  score: number;
  published_at: string;
  provider: string;
}

export interface ContentMetrics {
  views?: number;
  likes?: number;
  reading_time?: number;
  reactions?: number;
  recalculated_at?: string;
}

export interface ContentDetail extends Content {
  metrics?: ContentMetrics;
  updated_at?: string;
  full_content?: string;
}

export interface SearchFilters {
  q?: string;
  type?: ContentType | 'all';
  sort?: 'score_desc' | 'score_asc' | 'date_desc' | 'date_asc';
  page?: number;
  page_size?: number;
}

export interface PaginationInfo {
  page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export interface SearchResponse {
  success: boolean;
  data: Content[];
  pagination: PaginationInfo;
}

export interface ContentDetailResponse {
  success: boolean;
  data: ContentDetail;
}

export interface ProviderStats {
  provider_id: string;
  content_count: number;
  last_sync?: string;
}

export interface Stats {
  total_contents: number;
  total_videos: number;
  total_texts: number;
  average_score: number;
  last_sync?: string;
  providers: ProviderStats[];
}

export interface StatsResponse {
  success: boolean;
  data: Stats;
}

export interface ApiError {
  status: number;
  message: string;
  code?: string;
}


