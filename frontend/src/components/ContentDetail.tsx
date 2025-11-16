import { ArrowLeft, Star, Calendar, Eye, ThumbsUp, BookOpen, Heart, ExternalLink, Video, FileText, Clock } from 'lucide-react';
import { SearchResult } from '../types';
import { useMemo } from 'react';
import { useContentDetail } from '../hooks/useContentDetail';

interface ContentDetailProps {
  content: SearchResult;
  onBack: () => void;
}

export function ContentDetail({ content, onBack }: ContentDetailProps) {
  const { content: apiDetail, loading, error } = useContentDetail(String(content.id));

  const detail: SearchResult = useMemo(() => {
    if (!apiDetail) return content;
    return {
      id: String(apiDetail.id ?? content.id),
      title: apiDetail.title ?? content.title,
      description: apiDetail.description ?? content.description ?? '',
      contentType: apiDetail.content_type ?? content.contentType,
      score: apiDetail.score ?? content.score ?? 0,
      thumbnail: apiDetail.thumbnail_url ?? content.thumbnail,
      publishedDate: apiDetail.published_at ? new Date(apiDetail.published_at) : content.publishedDate,
      fullContent: apiDetail.full_content ?? content.fullContent,
      url: apiDetail.url ?? content.url,
      metrics: apiDetail.metrics ? {
        views: apiDetail.metrics.views,
        likes: apiDetail.metrics.likes,
        readingTime: apiDetail.metrics.reading_time,
        reactions: apiDetail.metrics.reactions,
      } : content.metrics,
      provider: apiDetail.provider ?? content.provider,
      lastUpdated: apiDetail.updated_at ? new Date(apiDetail.updated_at) : content.lastUpdated,
    };
  }, [apiDetail, content]);

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

  const formatNumber = (num: number): string => {
    if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
    return num.toString();
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-gray-400">Yükleniyor...</div>
    );
  }
  if (error) {
    // Hata olsa bile en azından listeden gelen içeriği gösterelim
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950 relative overflow-hidden">
      {/* Animated Background Elements */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl animate-pulse"></div>
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-green-500/10 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '1s' }}></div>
      </div>

      {/* Grid Pattern Overlay */}
      <div className="fixed inset-0 bg-[linear-gradient(rgba(16,185,129,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(16,185,129,0.03)_1px,transparent_1px)] bg-[size:50px_50px] [mask-image:radial-gradient(ellipse_80%_80%_at_50%_50%,black,transparent)]"></div>

      {/* Content */}
      <div className="relative z-10 max-w-5xl mx-auto px-6 py-8">
        {/* Back Button */}
        <button
          onClick={onBack}
          className="group flex items-center gap-2 mb-8 px-4 py-2 bg-gray-900/50 backdrop-blur-sm border border-gray-800 hover:border-emerald-500/50 rounded-lg text-gray-400 hover:text-emerald-400 transition-all duration-300"
        >
          <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform duration-300" />
          <span className="text-sm font-medium">Aramaya Dön</span>
        </button>

        {/* Main Content Card */}
        <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-2xl overflow-hidden">
          {/* Header Section */}
          <div className="p-8 border-b border-gray-800">
            <div className="flex items-start justify-between gap-6 mb-6">
              <div className="flex-1 space-y-4">
                <div className="flex items-center gap-3">
                  <div className={`px-3 py-1.5 ${detail.contentType === 'video' ? 'bg-emerald-500/90' : 'bg-teal-500/90'} backdrop-blur-sm rounded-lg flex items-center gap-1.5`}>
                    {detail.contentType === 'video' ? (
                      <Video className="w-3.5 h-3.5 text-gray-950" />
                    ) : (
                      <FileText className="w-3.5 h-3.5 text-gray-950" />
                    )}
                    <span className="text-xs font-semibold text-gray-950">
                      {detail.contentType === 'video' ? 'Video' : 'Metin'}
                    </span>
                  </div>
                  <div className="flex items-center gap-1.5 text-sm text-gray-500">
                    <Calendar className="w-4 h-4" />
                    <span>{getRelativeTime(detail.publishedDate)}</span>
                  </div>
                </div>

                <h1 className="text-4xl font-bold text-gray-100 leading-tight">
                  {detail.title}
                </h1>

                <p className="text-lg text-gray-400 leading-relaxed">
                  {detail.description}
                </p>
              </div>

              {/* Score Display */}
              <div className="flex flex-col items-center gap-2 p-6 bg-gradient-to-br from-emerald-500/10 to-teal-500/10 border border-emerald-500/20 rounded-xl">
                <Star className="w-8 h-8 text-yellow-500 fill-yellow-500" />
                <span className="text-4xl font-bold bg-gradient-to-r from-emerald-400 to-teal-400 bg-clip-text text-transparent">
                  {detail.score.toFixed(1)}
                </span>
                <span className="text-xs text-gray-500 uppercase tracking-wide">Skor</span>
              </div>
            </div>

            {/* Metrics Section */}
            {detail.metrics && (
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {detail.contentType === 'video' && (
                  <>
                    {detail.metrics.views !== undefined && (
                      <div className="flex items-center gap-3 p-4 bg-gray-800/50 rounded-xl">
                        <div className="p-2 bg-emerald-500/10 rounded-lg">
                          <Eye className="w-5 h-5 text-emerald-400" />
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">Görüntüleme</p>
                          <p className="text-lg font-bold text-gray-100">{formatNumber(detail.metrics.views)}</p>
                        </div>
                      </div>
                    )}
                    {detail.metrics.likes !== undefined && (
                      <div className="flex items-center gap-3 p-4 bg-gray-800/50 rounded-xl">
                        <div className="p-2 bg-emerald-500/10 rounded-lg">
                          <ThumbsUp className="w-5 h-5 text-emerald-400" />
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">Beğeni</p>
                          <p className="text-lg font-bold text-gray-100">{formatNumber(detail.metrics.likes)}</p>
                        </div>
                      </div>
                    )}
                  </>
                )}
                {detail.contentType === 'text' && (
                  <>
                    {detail.metrics.readingTime !== undefined && (
                      <div className="flex items-center gap-3 p-4 bg-gray-800/50 rounded-xl">
                        <div className="p-2 bg-teal-500/10 rounded-lg">
                          <BookOpen className="w-5 h-5 text-teal-400" />
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">Okuma Süresi</p>
                          <p className="text-lg font-bold text-gray-100">{detail.metrics.readingTime} dk</p>
                        </div>
                      </div>
                    )}
                    {detail.metrics.reactions !== undefined && (
                      <div className="flex items-center gap-3 p-4 bg-gray-800/50 rounded-xl">
                        <div className="p-2 bg-teal-500/10 rounded-lg">
                          <Heart className="w-5 h-5 text-teal-400" />
                        </div>
                        <div>
                          <p className="text-sm text-gray-500">Reaksiyon</p>
                          <p className="text-lg font-bold text-gray-100">{formatNumber(detail.metrics.reactions)}</p>
                        </div>
                      </div>
                    )}
                  </>
                )}
                {detail.lastUpdated && (
                  <div className="flex items-center gap-3 p-4 bg-gray-800/50 rounded-xl">
                    <div className="p-2 bg-gray-700/50 rounded-lg">
                      <Clock className="w-5 h-5 text-gray-400" />
                    </div>
                    <div>
                      <p className="text-sm text-gray-500">Son Güncelleme</p>
                      <p className="text-lg font-bold text-gray-100">{getRelativeTime(detail.lastUpdated)}</p>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Thumbnail/Image Section */}
          {detail.contentType === 'video' && detail.thumbnail && (
            <div className="relative h-96 overflow-hidden">
              <img
                src={detail.thumbnail}
                alt={detail.title}
                className="w-full h-full object-cover"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-gray-900 to-transparent"></div>
            </div>
          )}

          {/* Content Body */}
          <div className="p-8 space-y-6">
            {detail.fullContent && (
              <div className="prose prose-invert prose-emerald max-w-none">
                <p className="text-gray-300 leading-relaxed text-lg">
                  {detail.fullContent}
                </p>
              </div>
            )}

            {/* Action Buttons */}
            <div className="flex flex-wrap gap-4 pt-6 border-t border-gray-800">
              {detail.url && (
                <a
                  href={detail.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-emerald-500 to-teal-600 hover:from-emerald-400 hover:to-teal-500 text-gray-950 font-semibold rounded-lg transition-all duration-300 shadow-lg shadow-emerald-500/20 hover:shadow-emerald-500/30 hover:scale-105"
                >
                  <ExternalLink className="w-5 h-5" />
                  <span>Kaynağı Ziyaret Et</span>
                </a>
              )}
            </div>
          </div>

          {/* Provider Info */}
          {detail.provider && (
            <div className="p-6 bg-gray-800/30 border-t border-gray-800">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-500 mb-1">Sağlayıcı</p>
                  <p className="text-lg font-semibold text-gray-100">{detail.provider}</p>
                </div>
                <div className="px-4 py-2 bg-emerald-500/10 border border-emerald-500/20 rounded-lg">
                  <span className="text-sm font-medium text-emerald-400">Doğrulanmış</span>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
