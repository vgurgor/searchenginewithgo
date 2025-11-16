import { Filter, ChevronDown } from 'lucide-react';
import { ContentType, SortOption } from '../types';
import { useState } from 'react';

interface FilterBarProps {
  contentType: ContentType;
  sortBy: SortOption;
  onContentTypeChange: (type: ContentType) => void;
  onSortChange: (sort: SortOption) => void;
}

export function FilterBar({ contentType, sortBy, onContentTypeChange, onSortChange }: FilterBarProps) {
  const [showFilters, setShowFilters] = useState(false);

  const sortOptions: { value: SortOption; label: string }[] = [
    { value: 'score-high', label: 'Skor: Yüksekten Düşüğe' },
    { value: 'score-low', label: 'Skor: Düşükten Yükseğe' },
    { value: 'date-new', label: 'Tarih: En Yeni' },
    { value: 'date-old', label: 'Tarih: En Eski' },
  ];

  return (
    <div className="space-y-4">
      <button
        onClick={() => setShowFilters(!showFilters)}
        className="flex items-center gap-2 px-4 py-2 bg-gray-900/50 backdrop-blur-sm border border-gray-800 hover:border-emerald-500/50 rounded-lg text-gray-300 hover:text-emerald-400 transition-all duration-300"
      >
        <Filter className="w-4 h-4" />
        <span className="text-sm font-medium">Filtreler</span>
        <ChevronDown className={`w-4 h-4 transition-transform duration-300 ${showFilters ? 'rotate-180' : ''}`} />
      </button>

      {showFilters && (
        <div className="p-6 bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl space-y-6 animate-fade-in">
          <div className="space-y-3">
            <label className="text-sm font-medium text-gray-400">İçerik Tipi</label>
            <div className="flex gap-3">
              {(['all', 'video', 'text'] as ContentType[]).map((type) => (
                <button
                  key={type}
                  onClick={() => onContentTypeChange(type)}
                  className={`flex-1 px-4 py-2.5 rounded-lg font-medium text-sm transition-all duration-300 ${
                    contentType === type
                      ? 'bg-gradient-to-r from-emerald-500 to-teal-600 text-gray-950 shadow-lg shadow-emerald-500/25'
                      : 'bg-gray-800/50 text-gray-400 hover:text-gray-300 border border-gray-700 hover:border-gray-600'
                  }`}
                >
                  {type === 'all' ? 'Tümü' : type === 'video' ? 'Video' : 'Metin'}
                </button>
              ))}
            </div>
          </div>

          <div className="space-y-3">
            <label className="text-sm font-medium text-gray-400">Sıralama</label>
            <div className="relative">
              <select
                value={sortBy}
                onChange={(e) => onSortChange(e.target.value as SortOption)}
                className="w-full px-4 py-2.5 bg-gray-800/50 border border-gray-700 rounded-lg text-gray-300 outline-none cursor-pointer hover:border-gray-600 transition-colors duration-300 appearance-none"
              >
                {sortOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
              <ChevronDown className="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500 pointer-events-none" />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
