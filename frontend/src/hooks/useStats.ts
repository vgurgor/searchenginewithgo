import { useState, useEffect, useCallback } from 'react';
import { contentService } from '../api/content.service';
import { Stats, ApiError } from '../types/content.types';

interface UseStatsResult {
  stats: Stats | null;
  loading: boolean;
  error: ApiError | null;
  refetch: () => void;
}

export const useStats = (): UseStatsResult => {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<ApiError | null>(null);

  const fetchStats = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await contentService.getStats();
      setStats(response.data);
    } catch (err: any) {
      setError(err as ApiError);
      setStats(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  return {
    stats,
    loading,
    error,
    refetch: fetchStats,
  };
};


