import { useState, useEffect, useCallback } from 'react';
import { contentService } from '../api/content.service';
import {
  Content,
  SearchFilters,
  PaginationInfo,
  ApiError,
} from '../types/content.types';

interface UseSearchContentsResult {
  contents: Content[];
  pagination: PaginationInfo | null;
  loading: boolean;
  error: ApiError | null;
  searchContents: (filters: SearchFilters) => Promise<void>;
  refetch: () => void;
}

export const useSearchContents = (
  initialFilters: SearchFilters = {}
): UseSearchContentsResult => {
  const [contents, setContents] = useState<Content[]>([]);
  const [pagination, setPagination] = useState<PaginationInfo | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<ApiError | null>(null);
  const [currentFilters, setCurrentFilters] = useState<SearchFilters>(initialFilters);

  const searchContents = useCallback(async (filters: SearchFilters) => {
    setLoading(true);
    setError(null);
    setCurrentFilters(filters);
    try {
      const response = await contentService.searchContents(filters);
      setContents(response.data);
      setPagination(response.pagination);
    } catch (err: any) {
      setError(err as ApiError);
      setContents([]);
      setPagination(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const refetch = useCallback(() => {
    searchContents(currentFilters);
  }, [currentFilters, searchContents]);

  // İlk yüklemede otomatik arama yapma - App.tsx kontrol eder
  // useEffect(() => {
  //   searchContents(initialFilters);
  // }, []);

  return {
    contents,
    pagination,
    loading,
    error,
    searchContents,
    refetch,
  };
};


