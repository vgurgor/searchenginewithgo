import { SearchResult, ContentType, SortOption } from '../types';

export function filterAndSortResults(
  results: SearchResult[],
  query: string,
  contentType: ContentType,
  sortBy: SortOption
): SearchResult[] {
  let filtered = [...results];

  if (query.trim()) {
    const lowerQuery = query.toLowerCase();
    filtered = filtered.filter(
      (result) =>
        result.title.toLowerCase().includes(lowerQuery) ||
        result.description.toLowerCase().includes(lowerQuery)
    );
  }

  if (contentType !== 'all') {
    filtered = filtered.filter((result) => result.contentType === contentType);
  }

  filtered.sort((a, b) => {
    switch (sortBy) {
      case 'score-high':
        return b.score - a.score;
      case 'score-low':
        return a.score - b.score;
      case 'date-new':
        return b.publishedDate.getTime() - a.publishedDate.getTime();
      case 'date-old':
        return a.publishedDate.getTime() - b.publishedDate.getTime();
      default:
        return 0;
    }
  });

  return filtered;
}

export function paginateResults(
  results: SearchResult[],
  page: number,
  perPage: number
): SearchResult[] {
  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;
  return results.slice(startIndex, endIndex);
}
