import { useState, useEffect } from 'react';
import { contentService } from '../api/content.service';
import { ContentDetail, ApiError } from '../types/content.types';

interface UseContentDetailResult {
  content: ContentDetail | null;
  loading: boolean;
  error: ApiError | null;
}

export const useContentDetail = (id: string): UseContentDetailResult => {
  const [content, setContent] = useState<ContentDetail | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<ApiError | null>(null);

  useEffect(() => {
    const fetchContent = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await contentService.getContentById(id);
        setContent(response.data);
      } catch (err: any) {
        setError(err as ApiError);
        setContent(null);
      } finally {
        setLoading(false);
      }
    };
    if (id) {
      fetchContent();
    }
  }, [id]);

  return { content, loading, error };
};


