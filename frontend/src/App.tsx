import { useState, useEffect } from 'react';
import { Sparkles, BarChart3 } from 'lucide-react';
import { SearchBar } from './components/SearchBar';
import { FilterBar } from './components/FilterBar';
import { ResultCard } from './components/ResultCard';
import { Pagination } from './components/Pagination';
import { ContentDetail } from './components/ContentDetail';
import { StatsDashboard } from './components/StatsDashboard';
import { SkeletonCard } from './components/ResultCard';
import { useSearchContents } from './hooks/useSearchContents';
import { SearchFilters as ApiSearchFilters, Content as ApiContent } from './types/content.types';
import { ContentType, SortOption, SearchResult } from './types';

type View = 'search' | 'stats';

function App() {
  const [currentView, setCurrentView] = useState<View>('search');
  const [searchQuery, setSearchQuery] = useState('');
  const [contentType, setContentType] = useState<ContentType>('all');
  const [sortBy, setSortBy] = useState<SortOption>('score-high');
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedContent, setSelectedContent] = useState<SearchResult | null>(null);
  const [hasSearched, setHasSearched] = useState(false); // Arama yapıldı mı kontrolü
  const resultsPerPage = 9;
  const [results, setResults] = useState<SearchResult[]>([]);
  const [totalItems, setTotalItems] = useState(0);
  const totalPages = Math.ceil(totalItems / resultsPerPage);

  const { contents, pagination, loading, error, searchContents } = useSearchContents();

  // Sort eşleştirmeleri: UI <-> API
  const sortToApiMap: Record<SortOption, 'score_desc'|'score_asc'|'date_desc'|'date_asc'> = {
    'score-high': 'score_desc',
    'score-low': 'score_asc',
    'date-new': 'date_desc',
    'date-old': 'date_asc',
  };
  const apiToSortMap: Record<string, SortOption> = {
    'score_desc': 'score-high',
    'score_asc': 'score-low',
    'date_desc': 'date-new',
    'date_asc': 'date-old',
  };

  // İlk yüklemede URL query parametrelerini state'e uygula
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const q = params.get('q');
    const type = params.get('type');
    const sort = params.get('sort');
    const page = params.get('page');

    // URL'de parametre varsa state'leri doldur ve aramayı aktif et
    let hasUrlParams = false;
    
    if (q) {
      setSearchQuery(q);
      hasUrlParams = true;
    }
    if (type === 'video' || type === 'text') {
      setContentType(type as ContentType);
      hasUrlParams = true;
    }
    if (sort && (sort in apiToSortMap)) {
      setSortBy(apiToSortMap[sort as keyof typeof apiToSortMap]);
    }
    if (page) {
      const p = parseInt(page, 10);
      if (!Number.isNaN(p) && p > 0) setCurrentPage(p);
    }
    
    // URL'de arama parametresi varsa arama yapılmış sayalım
    if (hasUrlParams) {
      setHasSearched(true);
    }
  }, []);

  useEffect(() => {
    // Sadece arama yapılmışsa API'ye istek at
    if (!hasSearched) return;
    
    const typeParam = contentType === 'all' ? undefined : contentType;
    const filters: ApiSearchFilters = {
      q: searchQuery || undefined,
      type: typeParam as 'video' | 'text' | undefined,
      sort: sortToApiMap[sortBy as keyof typeof sortToApiMap],
      page: currentPage,
      page_size: resultsPerPage,
    };
    searchContents(filters);
  }, [searchQuery, contentType, sortBy, currentPage, hasSearched]);

  // Arama parametreleri değiştikçe URL'i güncelle (sadece arama yapılmışsa)
  useEffect(() => {
    if (!hasSearched) return;
    
    const params = new URLSearchParams();
    if (searchQuery) params.set('q', searchQuery);
    if (contentType !== 'all') params.set('type', contentType);
    params.set('sort', sortToApiMap[sortBy as keyof typeof sortToApiMap]);
    params.set('page', String(currentPage));
    params.set('page_size', String(resultsPerPage));
    const newUrl = `${window.location.pathname}?${params.toString()}`;
    window.history.replaceState(null, '', newUrl);
  }, [searchQuery, contentType, sortBy, currentPage, hasSearched]);

  useEffect(() => {
    const mapped: SearchResult[] = (contents || []).map((item: ApiContent) => ({
      id: String(item.id),
      title: item.title,
      description: item.description || '',
      contentType: item.content_type,
      score: item.score || 0,
      thumbnail: item.thumbnail_url || undefined,
      publishedDate: item.published_at ? new Date(item.published_at) : new Date(),
      url: item.url || undefined,
      provider: item.provider || undefined,
    }));
    setResults(mapped);
    setTotalItems(pagination?.total_items ? Number(pagination.total_items) : mapped.length);
  }, [contents, pagination]);

  useEffect(() => {
    if (hasSearched) {
      setCurrentPage(1);
    }
  }, [searchQuery, contentType, sortBy, hasSearched]);

  const handleSearch = () => {
    setHasSearched(true); // Arama yapıldığını işaretle
    setCurrentPage(1);
  };

  const handleContentClick = (content: SearchResult) => {
    setSelectedContent(content);
  };

  const handleBackToSearch = () => {
    setSelectedContent(null);
  };

  if (selectedContent) {
    return <ContentDetail content={selectedContent} onBack={handleBackToSearch} />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950 relative overflow-hidden">
      {/* Animated Background Elements */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl animate-pulse"></div>
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-green-500/10 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '1s' }}></div>
        <div className="absolute top-1/2 left-1/2 w-96 h-96 bg-teal-500/5 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '2s' }}></div>
      </div>

      {/* Grid Pattern Overlay */}
      <div className="fixed inset-0 bg-[linear-gradient(rgba(16,185,129,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(16,185,129,0.03)_1px,transparent_1px)] bg-[size:50px_50px] [mask-image:radial-gradient(ellipse_80%_80%_at_50%_50%,black,transparent)]"></div>

      {/* Header */}
      <header className="relative z-10 bg-gray-950/50 backdrop-blur-xl border-b border-gray-800">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3 group cursor-pointer">
              <div className="p-2 bg-gradient-to-br from-emerald-500 to-teal-600 rounded-lg shadow-lg shadow-emerald-500/25 group-hover:shadow-emerald-500/40 transition-all duration-300">
                <Sparkles className="w-6 h-6 text-gray-950" />
              </div>
              <span className="text-2xl font-bold bg-gradient-to-r from-emerald-400 via-green-400 to-teal-400 bg-clip-text text-transparent">
                Search Engine with Go
              </span>
            </div>
            <nav className="flex gap-3">
              <button
                onClick={() => { setCurrentView('search'); setSelectedContent(null); }}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-300 ${
                  currentView === 'search'
                    ? 'bg-gradient-to-r from-emerald-500 to-teal-600 text-gray-950 shadow-lg shadow-emerald-500/25'
                    : 'text-gray-400 hover:text-emerald-400 bg-gray-800/50'
                }`}
              >
                Arama
              </button>
              <button
                onClick={() => { setCurrentView('stats'); setSelectedContent(null); }}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-300 ${
                  currentView === 'stats'
                    ? 'bg-gradient-to-r from-emerald-500 to-teal-600 text-gray-950 shadow-lg shadow-emerald-500/25'
                    : 'text-gray-400 hover:text-emerald-400 bg-gray-800/50'
                }`}
              >
                <BarChart3 className="w-4 h-4" />
                <span>İstatistikler</span>
              </button>
            </nav>
          </div>

          {currentView === 'search' && (
            <div className="mt-6">
              <SearchBar value={searchQuery} onChange={setSearchQuery} onSearch={handleSearch} />
            </div>
          )}
        </div>
      </header>

      {/* Main Content */}
      <main className="relative z-10 max-w-7xl mx-auto px-6 py-8">
        {currentView === 'search' ? (
          <div className="space-y-6">
            <FilterBar
              contentType={contentType}
              sortBy={sortBy}
              onContentTypeChange={(type) => {
                setContentType(type);
                setHasSearched(true); // Filtre değiştiğinde arama yap
              }}
              onSortChange={(sort) => {
                setSortBy(sort);
                setHasSearched(true); // Sıralama değiştiğinde arama yap
              }}
            />

            {!hasSearched ? (
              <div className="text-center py-32">
                <div className="inline-block p-12 bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-2xl space-y-6">
                  <div className="flex justify-center">
                    <div className="p-6 bg-gradient-to-br from-emerald-500/10 to-teal-500/10 rounded-full">
                      <Sparkles className="w-16 h-16 text-emerald-400" />
                    </div>
                  </div>
                  <div className="space-y-3">
                    <h2 className="text-3xl font-bold text-gray-100">Aramaya Başlayın</h2>
                    <p className="text-gray-400 text-lg max-w-md">
                      Arama çubuğunu kullanarak içerikleri keşfedin veya filtreleri kullanarak sonuçları özelleştirin
                    </p>
                  </div>
                  <div className="flex flex-wrap justify-center gap-2 pt-4">
                    <span className="px-3 py-1.5 bg-gray-800/50 border border-gray-700 rounded-lg text-sm text-gray-400">
                      🎥 7 Video
                    </span>
                    <span className="px-3 py-1.5 bg-gray-800/50 border border-gray-700 rounded-lg text-sm text-gray-400">
                      📄 1 Metin
                    </span>
                    <span className="px-3 py-1.5 bg-gray-800/50 border border-gray-700 rounded-lg text-sm text-gray-400">
                      ⭐ Sıralama
                    </span>
                  </div>
                </div>
              </div>
            ) : loading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {Array.from({ length: 9 }, (_, i) => (
                  <div key={i}>
                    <SkeletonCard
                      type={i % 2 === 0 ? 'video' : 'text'}
                    />
                  </div>
                ))}
              </div>
            ) : error ? (
              <div className="text-center py-20">
                <div className="inline-block p-6 bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-2xl">
                  <p className="text-red-400 text-lg">Hata: {error.message}</p>
                </div>
              </div>
            ) : results.length === 0 ? (
              <div className="text-center py-20">
                <div className="inline-block p-6 bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-2xl">
                  <p className="text-gray-400 text-lg">Sonuç bulunamadı</p>
                  <p className="text-gray-500 text-sm mt-2">Farklı bir arama terimi veya filtre deneyin</p>
                </div>
              </div>
            ) : (
              <>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {results.map((result: SearchResult) => (
                    <ResultCard key={result.id} result={result} onClick={() => handleContentClick(result)} />
                  ))}
                </div>

                {totalPages > 1 && (
                  <Pagination
                    currentPage={currentPage}
                    totalPages={totalPages}
                    totalResults={totalItems}
                    resultsPerPage={resultsPerPage}
                    onPageChange={setCurrentPage}
                  />
                )}
              </>
            )}
          </div>
        ) : (
          <StatsDashboard />
        )}
      </main>

      {/* Footer */}
      <footer className="relative z-10 py-8 mt-20 border-t border-gray-800">
        <div className="max-w-7xl mx-auto px-6 text-center">
          <p className="text-gray-600 text-sm">
            &copy; 2025 Search Engine with Go. Yeni nesil arama deneyimi.
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
