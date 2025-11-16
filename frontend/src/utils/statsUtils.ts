import { SearchResult } from '../types';

export interface Stats {
  totalContents: number;
  totalVideos: number;
  totalTexts: number;
  averageScore: number;
}

export interface ProviderStats {
  name: string;
  contentCount: number;
  percentage: number;
  lastSync: Date;
}

export interface ContentDistribution {
  video: number;
  text: number;
}

export interface ScoreDistribution {
  range: string;
  count: number;
}

export function calculateStats(results: SearchResult[]): Stats {
  const totalContents = results.length;
  const totalVideos = results.filter(r => r.contentType === 'video').length;
  const totalTexts = results.filter(r => r.contentType === 'text').length;
  const averageScore = results.reduce((sum, r) => sum + r.score, 0) / totalContents;

  return {
    totalContents,
    totalVideos,
    totalTexts,
    averageScore: parseFloat(averageScore.toFixed(1)),
  };
}

export function calculateProviderStats(results: SearchResult[]): ProviderStats[] {
  const providerMap = new Map<string, { count: number; lastUpdate: Date }>();

  results.forEach(result => {
    if (result.provider) {
      const existing = providerMap.get(result.provider);
      if (existing) {
        existing.count++;
        if (result.lastUpdated && result.lastUpdated > existing.lastUpdate) {
          existing.lastUpdate = result.lastUpdated;
        }
      } else {
        providerMap.set(result.provider, {
          count: 1,
          lastUpdate: result.lastUpdated || result.publishedDate,
        });
      }
    }
  });

  const totalContents = results.length;
  const providerStats: ProviderStats[] = [];

  providerMap.forEach((value, key) => {
    providerStats.push({
      name: key,
      contentCount: value.count,
      percentage: (value.count / totalContents) * 100,
      lastSync: value.lastUpdate,
    });
  });

  return providerStats.sort((a, b) => b.contentCount - a.contentCount);
}

export function calculateContentDistribution(results: SearchResult[]): ContentDistribution {
  const videoCount = results.filter(r => r.contentType === 'video').length;
  const textCount = results.filter(r => r.contentType === 'text').length;

  return {
    video: videoCount,
    text: textCount,
  };
}

export function calculateScoreDistribution(results: SearchResult[]): ScoreDistribution[] {
  const ranges = [
    { range: '9.0 - 10.0', min: 9.0, max: 10.0, count: 0 },
    { range: '8.0 - 8.9', min: 8.0, max: 8.9, count: 0 },
    { range: '7.0 - 7.9', min: 7.0, max: 7.9, count: 0 },
    { range: '6.0 - 6.9', min: 6.0, max: 6.9, count: 0 },
    { range: '0.0 - 5.9', min: 0.0, max: 5.9, count: 0 },
  ];

  results.forEach(result => {
    const range = ranges.find(r => result.score >= r.min && result.score <= r.max);
    if (range) {
      range.count++;
    }
  });

  return ranges.map(r => ({ range: r.range, count: r.count }));
}
