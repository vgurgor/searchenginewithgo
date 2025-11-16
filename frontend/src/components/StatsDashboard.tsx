import { FileText, Video, BarChart3, Star, TrendingUp, CheckCircle } from 'lucide-react';
import { useStats } from '../hooks/useStats';

export function StatsDashboard() {
  const { stats, loading, error } = useStats();

  const contentDistribution = {
    video: stats?.total_videos || 0,
    text: stats?.total_texts || 0,
  };

  const overviewCards = [
    {
      title: 'Toplam İçerik',
      value: stats?.total_contents || 0,
      icon: BarChart3,
      color: 'from-emerald-500 to-teal-600',
      bgColor: 'from-emerald-500/10 to-teal-500/10',
      borderColor: 'border-emerald-500/20',
    },
    {
      title: 'Video İçerikler',
      value: stats?.total_videos || 0,
      icon: Video,
      color: 'from-green-500 to-emerald-600',
      bgColor: 'from-green-500/10 to-emerald-500/10',
      borderColor: 'border-green-500/20',
    },
    {
      title: 'Metin İçerikler',
      value: stats?.total_texts || 0,
      icon: FileText,
      color: 'from-teal-500 to-cyan-600',
      bgColor: 'from-teal-500/10 to-cyan-500/10',
      borderColor: 'border-teal-500/20',
    },
    {
      title: 'Ortalama Skor',
      value: (stats?.average_score || 0).toFixed(1),
      icon: Star,
      color: 'from-yellow-500 to-orange-600',
      bgColor: 'from-yellow-500/10 to-orange-500/10',
      borderColor: 'border-yellow-500/20',
    },
  ];

  if (loading) {
    return <div className="text-gray-400">Yükleniyor...</div>;
  }
  if (error) {
    return <div className="text-red-400">Hata: {error.message}</div>;
  }

  return (
    <div className="space-y-8">
      {/* Overview Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {overviewCards.map((card, index) => {
          const Icon = card.icon;
          return (
            <div
              key={index}
              className={`relative group bg-gradient-to-br ${card.bgColor} backdrop-blur-sm border ${card.borderColor} rounded-xl p-6 overflow-hidden transition-all duration-300 hover:scale-105`}
            >
              <div className={`absolute top-0 right-0 w-32 h-32 bg-gradient-to-br ${card.color} opacity-10 rounded-full blur-2xl group-hover:opacity-20 transition-opacity duration-300`}></div>

              <div className="relative flex items-center justify-between">
                <div className="space-y-2">
                  <p className="text-sm text-gray-400 font-medium">{card.title}</p>
                  <p className="text-4xl font-bold text-gray-100">{card.value}</p>
                </div>
                <div className={`p-4 bg-gradient-to-br ${card.color} rounded-xl shadow-lg`}>
                  <Icon className="w-8 h-8 text-gray-950" />
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 gap-6">
        {/* Content Distribution */}
        <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl p-6">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-2 bg-emerald-500/10 rounded-lg">
              <TrendingUp className="w-5 h-5 text-emerald-400" />
            </div>
            <h2 className="text-xl font-bold text-gray-100">İçerik Dağılımı</h2>
          </div>

          <div className="space-y-6">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Video className="w-4 h-4 text-emerald-400" />
                  <span className="text-gray-300 font-medium">Video</span>
                </div>
                <span className="text-gray-100 font-bold">{contentDistribution.video}</span>
              </div>
              <div className="relative h-3 bg-gray-800 rounded-full overflow-hidden">
                <div
                  className="absolute inset-y-0 left-0 bg-gradient-to-r from-emerald-500 to-teal-600 rounded-full transition-all duration-500"
                  style={{ width: `${(stats?.total_contents ? (contentDistribution.video / stats.total_contents) * 100 : 0)}%` }}
                ></div>
              </div>
              <p className="text-sm text-gray-500 text-right">
                {(stats?.total_contents ? ((contentDistribution.video / stats.total_contents) * 100).toFixed(1) : '0.0')}%
              </p>
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <FileText className="w-4 h-4 text-teal-400" />
                  <span className="text-gray-300 font-medium">Metin</span>
                </div>
                <span className="text-gray-100 font-bold">{contentDistribution.text}</span>
              </div>
              <div className="relative h-3 bg-gray-800 rounded-full overflow-hidden">
                <div
                  className="absolute inset-y-0 left-0 bg-gradient-to-r from-teal-500 to-cyan-600 rounded-full transition-all duration-500"
                  style={{ width: `${(stats?.total_contents ? (contentDistribution.text / stats.total_contents) * 100 : 0)}%` }}
                ></div>
              </div>
              <p className="text-sm text-gray-500 text-right">
                {(stats?.total_contents ? ((contentDistribution.text / stats.total_contents) * 100).toFixed(1) : '0.0')}%
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Provider Statistics */}
      <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl overflow-hidden">
        <div className="p-6 border-b border-gray-800">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-emerald-500/10 rounded-lg">
              <CheckCircle className="w-5 h-5 text-emerald-400" />
            </div>
            <h2 className="text-xl font-bold text-gray-100">Sağlayıcı İstatistikleri</h2>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-800/50">
              <tr>
                <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">Sağlayıcı</th>
                <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">İçerik Sayısı</th>
                <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">Dağılım</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-800">
              {(stats?.providers || []).map((provider, index) => (
                <tr key={index} className="hover:bg-gray-800/30 transition-colors duration-200">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                      <span className="font-medium text-gray-100">{provider.provider_id}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <span className="text-gray-300 font-semibold">{provider.content_count}</span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="space-y-2">
                      <div className="relative h-2 bg-gray-800 rounded-full overflow-hidden w-48">
                        <div
                          className="absolute inset-y-0 left-0 bg-gradient-to-r from-emerald-500 to-teal-600 rounded-full transition-all duration-500"
                          style={{ width: `${stats?.total_contents ? (provider.content_count / stats.total_contents) * 100 : 0}%` }}
                        ></div>
                      </div>
                      <span className="text-xs text-gray-500">
                        {stats?.total_contents ? ((provider.content_count / stats.total_contents) * 100).toFixed(1) : 0}%
                      </span>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
