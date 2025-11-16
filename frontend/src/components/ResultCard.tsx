import { Video, FileText, Star, Calendar, Eye } from 'lucide-react';
import { SearchResult } from '../types';

interface ResultCardProps {
  result: SearchResult;
  onClick?: () => void;
}

interface SkeletonCardProps {
  type?: 'video' | 'text';
}

export function ResultCard({ result, onClick }: ResultCardProps) {
  const getRelativeTime = (date: Date): string => {
    const now = new Date();
    const diffInMs = now.getTime() - date.getTime();
    const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

    if (diffInDays === 0) return 'Bugün';
    if (diffInDays === 1) return 'Dün';
    if (diffInDays < 7) return `${diffInDays} gün önce`;
    if (diffInDays < 30) return `${Math.floor(diffInDays / 7)} hafta önce`;
    return `${Math.floor(diffInDays / 30)} ay önce`;
  };

  return (
    <div
      onClick={onClick}
      className="group relative bg-gray-900/50 backdrop-blur-sm border border-gray-800 hover:border-gray-700 rounded-xl overflow-hidden transition-all duration-300 hover:scale-[1.02] hover:shadow-xl hover:shadow-emerald-500/10 cursor-pointer"
    >
      <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/0 to-teal-500/0 group-hover:from-emerald-500/5 group-hover:to-teal-500/5 transition-all duration-300"></div>

      {result.contentType === 'video' && (
        <div className="relative h-48 overflow-hidden">
          {result.thumbnail ? (
            <img
              src={result.thumbnail}
              alt={result.title}
              className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
              onError={(e) => {
                const target = e.target as HTMLImageElement;
                target.style.display = 'none';
                const parent = target.parentElement;
                if (parent) {
                  parent.classList.add('bg-gradient-to-br', 'from-gray-800', 'to-gray-900', 'flex', 'items-center', 'justify-center');
                  const icon = document.createElement('div');
                  icon.innerHTML = '<svg class="w-16 h-16 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>';
                  parent.insertBefore(icon, target);
                }
              }}
            />
          ) : (
            <div className="w-full h-full bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center">
              <Video className="w-16 h-16 text-gray-700" />
            </div>
          )}
          <div className="absolute inset-0 bg-gradient-to-t from-gray-900 to-transparent"></div>
          <div className="absolute top-3 right-3 px-3 py-1.5 bg-emerald-500/90 backdrop-blur-sm rounded-lg flex items-center gap-1.5">
            <Video className="w-3.5 h-3.5 text-gray-950" />
            <span className="text-xs font-semibold text-gray-950">Video</span>
          </div>
        </div>
      )}

      {result.contentType === 'text' && (
        <div className="relative h-48 bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center">
          <FileText className="w-16 h-16 text-gray-700" />
          <div className="absolute top-3 right-3 px-3 py-1.5 bg-teal-500/90 backdrop-blur-sm rounded-lg flex items-center gap-1.5">
            <FileText className="w-3.5 h-3.5 text-gray-950" />
            <span className="text-xs font-semibold text-gray-950">Metin</span>
          </div>
        </div>
      )}

      <div className="relative p-5 space-y-3">
        <div className="flex items-start justify-between gap-3">
          <h3 className="flex-1 text-lg font-semibold text-gray-100 line-clamp-2 group-hover:text-emerald-400 transition-colors duration-300">
            {result.title}
          </h3>
          <div className="flex items-center gap-1.5 px-2.5 py-1 bg-gray-800/50 rounded-lg shrink-0">
            <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
            <span className="text-sm font-bold text-gray-100">{result.score.toFixed(1)}</span>
          </div>
        </div>

        <p className="text-sm text-gray-400 line-clamp-3 leading-relaxed">
          {result.description}
        </p>

        <div className="flex items-center justify-between pt-3 border-t border-gray-800">
          <div className="flex items-center gap-1.5 text-xs text-gray-500">
            <Calendar className="w-3.5 h-3.5" />
            <span>{getRelativeTime(result.publishedDate)}</span>
          </div>
          <div className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-emerald-500/10 to-teal-500/10 group-hover:from-emerald-500/20 group-hover:to-teal-500/20 border border-emerald-500/20 rounded-lg text-emerald-400 text-sm font-medium transition-all duration-300 group-hover:border-emerald-500/40">
            <Eye className="w-4 h-4" />
            <span>Detayları Gör</span>
          </div>
        </div>
      </div>
    </div>
  );
}

export function SkeletonCard({ type = 'video' }: SkeletonCardProps) {
  return (
    <div className="group relative bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/0 to-teal-500/0"></div>

      {type === 'video' && (
        <div className="relative h-48 bg-gray-800 animate-pulse">
          <div className="absolute inset-0 bg-gradient-to-t from-gray-900 to-transparent"></div>
          <div className="absolute top-3 right-3 px-3 py-1.5 bg-gray-700/90 backdrop-blur-sm rounded-lg flex items-center gap-1.5">
            <Video className="w-3.5 h-3.5 text-gray-600" />
            <span className="text-xs font-semibold text-gray-600">Video</span>
          </div>
        </div>
      )}

      {type === 'text' && (
        <div className="relative h-48 bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center animate-pulse">
          <FileText className="w-16 h-16 text-gray-700" />
          <div className="absolute top-3 right-3 px-3 py-1.5 bg-gray-700/90 backdrop-blur-sm rounded-lg flex items-center gap-1.5">
            <FileText className="w-3.5 h-3.5 text-gray-600" />
            <span className="text-xs font-semibold text-gray-600">Metin</span>
          </div>
        </div>
      )}

      <div className="relative p-5 space-y-3">
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1 space-y-2">
            <div className="h-5 bg-gray-800 rounded animate-pulse"></div>
            <div className="h-4 bg-gray-800 rounded w-3/4 animate-pulse"></div>
          </div>
          <div className="flex items-center gap-1.5 px-2.5 py-1 bg-gray-800/50 rounded-lg shrink-0">
            <Star className="w-4 h-4 text-gray-600" />
            <div className="w-8 h-4 bg-gray-800 rounded animate-pulse"></div>
          </div>
        </div>

        <div className="space-y-2">
          <div className="h-4 bg-gray-800 rounded animate-pulse"></div>
          <div className="h-4 bg-gray-800 rounded w-5/6 animate-pulse"></div>
          <div className="h-4 bg-gray-800 rounded w-4/6 animate-pulse"></div>
        </div>

        <div className="flex items-center justify-between pt-3 border-t border-gray-800">
          <div className="flex items-center gap-1.5 text-xs">
            <Calendar className="w-3.5 h-3.5 text-gray-600" />
            <div className="w-16 h-3 bg-gray-800 rounded animate-pulse"></div>
          </div>
          <div className="flex items-center gap-2 px-4 py-2 bg-gray-800/50 border border-gray-700 rounded-lg">
            <Eye className="w-4 h-4 text-gray-600" />
            <div className="w-20 h-4 bg-gray-800 rounded animate-pulse"></div>
          </div>
        </div>
      </div>
    </div>
  );
}
