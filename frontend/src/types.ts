export type ContentType = 'all' | 'video' | 'text';
export type SortOption = 'score-high' | 'score-low' | 'date-new' | 'date-old';

export interface SearchResult {
  id: string;
  title: string;
  description: string;
  contentType: 'video' | 'text';
  score: number;
  thumbnail?: string;
  publishedDate: Date;
  fullContent?: string;
  url?: string;
  metrics?: {
    views?: number;
    likes?: number;
    readingTime?: number;
    reactions?: number;
  };
  provider?: string;
  lastUpdated?: Date;
}

export interface SearchFilters {
  contentType: ContentType;
  sortBy: SortOption;
}
