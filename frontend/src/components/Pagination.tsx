import { ChevronLeft, ChevronRight } from 'lucide-react';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  totalResults: number;
  resultsPerPage: number;
  onPageChange: (page: number) => void;
}

export function Pagination({
  currentPage,
  totalPages,
  totalResults,
  resultsPerPage,
  onPageChange,
}: PaginationProps) {
  const startResult = (currentPage - 1) * resultsPerPage + 1;
  const endResult = Math.min(currentPage * resultsPerPage, totalResults);

  const getPageNumbers = () => {
    const pages: (number | string)[] = [];
    const showPages = 5;

    if (totalPages <= showPages + 2) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);

      if (currentPage > 3) {
        pages.push('...');
      }

      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (currentPage < totalPages - 2) {
        pages.push('...');
      }

      pages.push(totalPages);
    }

    return pages;
  };

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 pt-8">
      <p className="text-sm text-gray-400">
        <span className="font-medium text-gray-300">{startResult}-{endResult}</span> arası gösteriliyor,{' '}
        <span className="font-medium text-gray-300">{totalResults}</span> sonuç
      </p>

      <div className="flex items-center gap-2">
        <button
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className="p-2 bg-gray-900/50 backdrop-blur-sm border border-gray-800 hover:border-gray-700 rounded-lg text-gray-400 hover:text-emerald-400 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:text-gray-400"
        >
          <ChevronLeft className="w-5 h-5" />
        </button>

        <div className="flex items-center gap-1">
          {getPageNumbers().map((page, index) => (
            <button
              key={index}
              onClick={() => typeof page === 'number' && onPageChange(page)}
              disabled={page === '...'}
              className={`min-w-[40px] h-10 px-3 rounded-lg font-medium text-sm transition-all duration-300 ${
                page === currentPage
                  ? 'bg-gradient-to-r from-emerald-500 to-teal-600 text-gray-950 shadow-lg shadow-emerald-500/25'
                  : page === '...'
                  ? 'text-gray-600 cursor-default'
                  : 'bg-gray-900/50 backdrop-blur-sm border border-gray-800 hover:border-gray-700 text-gray-400 hover:text-emerald-400'
              }`}
            >
              {page}
            </button>
          ))}
        </div>

        <button
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="p-2 bg-gray-900/50 backdrop-blur-sm border border-gray-800 hover:border-gray-700 rounded-lg text-gray-400 hover:text-emerald-400 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:text-gray-400"
        >
          <ChevronRight className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
}
