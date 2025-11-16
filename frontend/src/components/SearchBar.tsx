import { Search, History, TrendingUp, X } from 'lucide-react';
import { useState, useEffect, useRef } from 'react';

interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  onSearch: () => void;
}

const POPULAR_SEARCHES = [
  'Go programming',
  'Python tutorial',
  'Web development',
  'Machine learning',
  'Docker containers',
  'React components',
  'Database design',
  'API development'
];

const MAX_RECENT_SEARCHES = 5;

export function SearchBar({ value, onChange, onSearch }: SearchBarProps) {
  const [isFocused, setIsFocused] = useState(false);
  const [recentSearches, setRecentSearches] = useState<string[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);

  // Load recent searches from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('recentSearches');
    if (saved) {
      try {
        setRecentSearches(JSON.parse(saved));
      } catch (e) {
        console.warn('Failed to parse recent searches');
      }
    }
  }, []);

  // Save recent searches to localStorage
  const saveRecentSearch = (search: string) => {
    if (!search.trim()) return;

    const updated = [search, ...recentSearches.filter(s => s !== search)].slice(0, MAX_RECENT_SEARCHES);
    setRecentSearches(updated);
    localStorage.setItem('recentSearches', JSON.stringify(updated));
  };

  // Handle search submission
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (value.trim()) {
      saveRecentSearch(value.trim());
      onSearch();
      setShowSuggestions(false);
    }
  };

  // Handle suggestion click
  const handleSuggestionClick = (suggestion: string) => {
    onChange(suggestion);
    saveRecentSearch(suggestion);
    onSearch();
    setShowSuggestions(false);
    inputRef.current?.blur();
  };

  // Handle input focus
  const handleFocus = () => {
    setIsFocused(true);
    setShowSuggestions(true);
  };

  // Handle input blur
  const handleBlur = () => {
    // Delay to allow suggestion clicks
    setTimeout(() => {
      setIsFocused(false);
      setShowSuggestions(false);
    }, 150);
  };

  // Clear recent searches
  const clearRecentSearches = () => {
    setRecentSearches([]);
    localStorage.removeItem('recentSearches');
  };

  // Filter suggestions based on input
  const getSuggestions = () => {
    if (!value.trim()) {
      return {
        recent: recentSearches,
        popular: POPULAR_SEARCHES.slice(0, 4)
      };
    }

    const query = value.toLowerCase();
    const filteredRecent = recentSearches.filter(s =>
      s.toLowerCase().includes(query)
    );
    const filteredPopular = POPULAR_SEARCHES.filter(s =>
      s.toLowerCase().includes(query)
    );

    return {
      recent: filteredRecent,
      popular: filteredPopular
    };
  };

  const suggestions = getSuggestions();

  return (
    <div className="relative group">
      <form onSubmit={handleSubmit} className="relative">
        <div className="absolute -inset-0.5 bg-gradient-to-r from-emerald-500 via-green-500 to-teal-500 rounded-xl blur opacity-20 group-hover:opacity-30 transition-opacity duration-300"></div>
        <div className="relative flex items-center bg-gray-900/90 backdrop-blur-xl border border-gray-800 rounded-xl overflow-hidden">
          <div className="pl-5 pr-3">
            <Search className="w-5 h-5 text-gray-500 group-hover:text-emerald-500 transition-colors duration-300" />
          </div>
          <input
            ref={inputRef}
            type="text"
            value={value}
            onChange={(e) => onChange(e.target.value)}
            onFocus={handleFocus}
            onBlur={handleBlur}
            placeholder="Ne aramak istersiniz?"
            className="flex-1 bg-transparent py-4 text-gray-100 placeholder-gray-500 outline-none"
            autoComplete="off"
          />
          {value && (
            <button
              type="button"
              onClick={() => onChange('')}
              className="pr-3 text-gray-500 hover:text-gray-300 transition-colors"
            >
              <X className="w-4 h-4" />
            </button>
          )}
          <button
            type="submit"
            disabled={!value.trim()}
            className="m-1 px-6 py-3 bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-400 hover:to-teal-500 disabled:from-gray-600 disabled:to-gray-700 disabled:cursor-not-allowed text-gray-950 disabled:text-gray-400 font-semibold rounded-lg transition-all duration-300 shadow-lg shadow-emerald-500/20 hover:shadow-emerald-500/30"
          >
            Ara
          </button>
        </div>
      </form>

      {/* Suggestions Dropdown */}
      {showSuggestions && (
        <div
          ref={suggestionsRef}
          className="absolute top-full left-0 right-0 mt-2 bg-gray-900/95 backdrop-blur-xl border border-gray-800 rounded-xl shadow-2xl z-50 max-h-96 overflow-y-auto"
        >
          {/* Recent Searches */}
          {suggestions.recent.length > 0 && (
            <div className="p-3">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2 text-sm text-gray-400 font-medium">
                  <History className="w-4 h-4" />
                  Son Aramalar
                </div>
                <button
                  onClick={clearRecentSearches}
                  className="text-xs text-gray-500 hover:text-gray-300 transition-colors"
                >
                  Temizle
                </button>
              </div>
              {suggestions.recent.map((search, index) => (
                <button
                  key={`recent-${index}`}
                  onClick={() => handleSuggestionClick(search)}
                  className="w-full text-left px-3 py-2 text-sm text-gray-200 hover:bg-gray-800/50 rounded-lg transition-colors flex items-center gap-3"
                >
                  <History className="w-4 h-4 text-gray-500" />
                  {search}
                </button>
              ))}
            </div>
          )}

          {/* Popular Searches */}
          {suggestions.popular.length > 0 && (
            <div className="p-3 border-t border-gray-800">
              <div className="flex items-center gap-2 mb-2">
                <TrendingUp className="w-4 h-4 text-gray-400" />
                <span className="text-sm text-gray-400 font-medium">Popüler Aramalar</span>
              </div>
              {suggestions.popular.map((search, index) => (
                <button
                  key={`popular-${index}`}
                  onClick={() => handleSuggestionClick(search)}
                  className="w-full text-left px-3 py-2 text-sm text-gray-200 hover:bg-gray-800/50 rounded-lg transition-colors flex items-center gap-3"
                >
                  <TrendingUp className="w-4 h-4 text-gray-500" />
                  {search}
                </button>
              ))}
            </div>
          )}

          {/* No suggestions */}
          {suggestions.recent.length === 0 && suggestions.popular.length === 0 && value.trim() && (
            <div className="p-4 text-center text-gray-500">
              Arama önerisi bulunamadı
            </div>
          )}
        </div>
      )}
    </div>
  );
}
