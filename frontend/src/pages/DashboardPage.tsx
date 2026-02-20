import { useEffect, useState, useCallback, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { MousePointerClick, Activity, Users, Bot, Globe, Eye } from 'lucide-react'
import { useClickStats, useEventStats, getDateRange, type Period } from '@/hooks/useStats'
import { useRpcList } from '@/hooks/useRpc'
import { RPC } from '@/lib/constants'
import { mergeNameCounts } from '@/lib/utils'
import { StatCard, RankingList, MiniPieChart, TrendChart, HourlyChart } from '@/components/dashboard'
import type { Tracker, Campaign, Channel, Site, GroupCount } from '@/types'

type TrendView = 'clicks' | 'events' | 'both'

function groupCountToNameCount(data: GroupCount[] | undefined) {
  return (data ?? []).map((d) => ({ name: d.name, count: d.count }))
}

export function DashboardPage() {
  const { t } = useTranslation()
  const [period, setPeriod] = useState<Period>('7d')
  const [customStart, setCustomStart] = useState('')
  const [customEnd, setCustomEnd] = useState('')
  const [trendView, setTrendView] = useState<TrendView>('both')

  // Filters
  const [filterTracker, setFilterTracker] = useState('')
  const [filterCampaign, setFilterCampaign] = useState('')
  const [filterChannel, setFilterChannel] = useState('')
  const [filterSite, setFilterSite] = useState('')

  const clickStats = useClickStats()
  const eventStats = useEventStats()

  // Filter dropdown data
  const trackers = useRpcList<Tracker>(RPC.TRACKER_LIST)
  const campaigns = useRpcList<Campaign>(RPC.CAMPAIGN_LIST)
  const channels = useRpcList<Channel>(RPC.CHANNEL_LIST)
  const sites = useRpcList<Site>(RPC.SITE_LIST)

  useEffect(() => {
    trackers.fetch()
    campaigns.fetch()
    channels.fetch()
    sites.fetch()
  }, [trackers.fetch, campaigns.fetch, channels.fetch, sites.fetch])

  const fetchAll = useCallback(() => {
    let startDate: string, endDate: string
    if (period === 'custom') {
      if (!customStart || !customEnd) return
      startDate = customStart
      endDate = customEnd
    } else {
      const range = getDateRange(period)
      startDate = range.startDate
      endDate = range.endDate
    }
    clickStats.fetch({
      start_date: startDate,
      end_date: endDate,
      ...(filterTracker && { tracker_id: filterTracker }),
      ...(filterCampaign && { campaign_id: filterCampaign }),
      ...(filterChannel && { channel_id: filterChannel }),
    })
    eventStats.fetch({
      start_date: startDate,
      end_date: endDate,
      ...(filterSite && { site_id: filterSite }),
    })
  }, [period, customStart, customEnd, filterTracker, filterCampaign, filterChannel, filterSite, clickStats.fetch, eventStats.fetch])

  useEffect(() => {
    fetchAll()
  }, [fetchAll])

  const loading = clickStats.loading || eventStats.loading

  // Trend data
  const trendData = useMemo(() => {
    const clickMap = new Map<string, number>()
    const eventMap = new Map<string, number>()
    clickStats.data?.daily?.forEach((d) => clickMap.set(d.date, d.count))
    eventStats.data?.daily?.forEach((d) => eventMap.set(d.date, d.count))
    const allDates = new Set([...clickMap.keys(), ...eventMap.keys()])
    return Array.from(allDates)
      .sort()
      .map((date) => ({
        date,
        clicks: clickMap.get(date) ?? 0,
        events: eventMap.get(date) ?? 0,
      }))
  }, [clickStats.data?.daily, eventStats.data?.daily])

  // Bot trend data
  const botTrendData = useMemo(() => {
    const clickMap = new Map<string, number>()
    const eventMap = new Map<string, number>()
    clickStats.data?.bot_daily?.forEach((d) => clickMap.set(d.date, d.count))
    eventStats.data?.bot_daily?.forEach((d) => eventMap.set(d.date, d.count))
    const allDates = new Set([...clickMap.keys(), ...eventMap.keys()])
    return Array.from(allDates)
      .sort()
      .map((date) => ({
        date,
        clickBots: clickMap.get(date) ?? 0,
        eventBots: eventMap.get(date) ?? 0,
      }))
  }, [clickStats.data?.bot_daily, eventStats.data?.bot_daily])

  // Merged distributions
  const browsers = useMemo(() => mergeNameCounts(clickStats.data?.browsers, eventStats.data?.browsers), [clickStats.data?.browsers, eventStats.data?.browsers])
  const oses = useMemo(() => mergeNameCounts(clickStats.data?.oses, eventStats.data?.oses), [clickStats.data?.oses, eventStats.data?.oses])
  const languages = useMemo(() => mergeNameCounts(clickStats.data?.languages, eventStats.data?.languages), [clickStats.data?.languages, eventStats.data?.languages])
  const referrers = useMemo(() => mergeNameCounts(clickStats.data?.top_referrers, eventStats.data?.top_referrers), [clickStats.data?.top_referrers, eventStats.data?.top_referrers])

  const trendSeries = useMemo(() => {
    switch (trendView) {
      case 'clicks':
        return [{ key: 'clicks', color: '#6366f1', name: t('dashboard.clicks') }]
      case 'events':
        return [{ key: 'events', color: '#8b5cf6', name: t('dashboard.events') }]
      default:
        return [
          { key: 'clicks', color: '#6366f1', name: t('dashboard.clicks') },
          { key: 'events', color: '#8b5cf6', name: t('dashboard.events') },
        ]
    }
  }, [trendView, t])

  const periods: { key: Period; label: string }[] = [
    { key: 'today', label: t('dashboard.today') },
    { key: '7d', label: t('dashboard.last7d') },
    { key: '30d', label: t('dashboard.last30d') },
    { key: 'custom', label: t('dashboard.custom') },
  ]

  const trendViews: { key: TrendView; label: string }[] = [
    { key: 'clicks', label: t('dashboard.trendViewClicks') },
    { key: 'events', label: t('dashboard.trendViewEvents') },
    { key: 'both', label: t('dashboard.trendViewBoth') },
  ]

  const clickTotal = clickStats.data?.summary?.total ?? 0
  const eventTotal = eventStats.data?.summary?.total ?? 0

  return (
    <div>
      <h2 className="mb-4 text-2xl font-bold font-mono text-text">{t('dashboard.title')}</h2>

      {/* Period selector + Filters */}
      <div className="mb-4 flex flex-wrap items-center gap-2">
        {periods.map((p) => (
          <button
            key={p.key}
            onClick={() => setPeriod(p.key)}
            className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
              period === p.key
                ? 'bg-primary text-white'
                : 'bg-bg-card border border-border text-text hover:bg-primary/10'
            }`}
          >
            {p.label}
          </button>
        ))}
        {period === 'custom' && (
          <div className="flex items-center gap-2 ml-2">
            <input
              type="date"
              value={customStart}
              onChange={(e) => setCustomStart(e.target.value)}
              className="rounded-lg border border-border bg-bg-card px-2 py-1.5 text-sm text-text"
            />
            <span className="text-muted">-</span>
            <input
              type="date"
              value={customEnd}
              onChange={(e) => setCustomEnd(e.target.value)}
              className="rounded-lg border border-border bg-bg-card px-2 py-1.5 text-sm text-text"
            />
          </div>
        )}
      </div>

      {/* Filter dropdowns */}
      <div className="mb-6 flex flex-wrap items-center gap-2">
        <select
          value={filterTracker}
          onChange={(e) => setFilterTracker(e.target.value)}
          className="rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text"
        >
          <option value="">{t('dashboard.allTrackers')}</option>
          {trackers.data.map((tr) => (
            <option key={tr.id} value={tr.id}>{tr.name}</option>
          ))}
        </select>
        <select
          value={filterCampaign}
          onChange={(e) => setFilterCampaign(e.target.value)}
          className="rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text"
        >
          <option value="">{t('dashboard.allCampaigns')}</option>
          {campaigns.data.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>
        <select
          value={filterChannel}
          onChange={(e) => setFilterChannel(e.target.value)}
          className="rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text"
        >
          <option value="">{t('dashboard.allChannels')}</option>
          {channels.data.map((ch) => (
            <option key={ch.id} value={ch.id}>{ch.name}</option>
          ))}
        </select>
        <select
          value={filterSite}
          onChange={(e) => setFilterSite(e.target.value)}
          className="rounded-lg border border-border bg-bg-card px-3 py-1.5 text-sm text-text"
        >
          <option value="">{t('dashboard.allSites')}</option>
          {sites.data.map((s) => (
            <option key={s.id} value={s.id}>{s.name}</option>
          ))}
        </select>
      </div>

      {/* 6 Stat cards */}
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6 mb-6">
        <StatCard icon={MousePointerClick} label={t('dashboard.totalClicks')} value={clickTotal} loading={loading} />
        <StatCard icon={Activity} label={t('dashboard.totalEvents')} value={eventTotal} loading={loading} />
        <StatCard icon={Users} label={t('dashboard.uniqueVisitorsClick')} value={clickStats.data?.summary?.unique_visitors ?? 0} loading={loading} />
        <StatCard icon={Eye} label={t('dashboard.uniqueVisitorsEvent')} value={eventStats.data?.summary?.unique_visitors ?? 0} loading={loading} />
        <StatCard icon={Bot} label={t('dashboard.clickBotRate')} value={`${clickStats.data?.summary?.bot_rate ?? 0}%`} loading={loading} />
        <StatCard icon={Globe} label={t('dashboard.eventBotRate')} value={`${eventStats.data?.summary?.bot_rate ?? 0}%`} loading={loading} />
      </div>

      {/* Trend chart with toggle */}
      <div className="rounded-xl border border-border bg-bg-card p-4 mb-6">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-sm font-semibold text-text">{t('dashboard.trend')}</h3>
          <div className="flex gap-1">
            {trendViews.map((tv) => (
              <button
                key={tv.key}
                onClick={() => setTrendView(tv.key)}
                className={`rounded px-2.5 py-1 text-xs font-medium transition-colors ${
                  trendView === tv.key
                    ? 'bg-primary text-white'
                    : 'text-muted hover:bg-primary/10'
                }`}
              >
                {tv.label}
              </button>
            ))}
          </div>
        </div>
        <TrendChart
          title=""
          data={trendData}
          series={trendSeries}
          height={280}
          stacked={trendView === 'both'}
        />
      </div>

      {/* Rankings - 3 columns */}
      <div className="grid grid-cols-1 gap-3 lg:grid-cols-3 mb-6">
        <RankingList title={t('dashboard.topTrackers')} data={groupCountToNameCount(clickStats.data?.top_trackers)} />
        <RankingList title={t('dashboard.topCampaigns')} data={groupCountToNameCount(clickStats.data?.top_campaigns)} />
        <RankingList title={t('dashboard.topChannels')} data={groupCountToNameCount(clickStats.data?.top_channels)} />
        <RankingList title={t('dashboard.topSites')} data={groupCountToNameCount(eventStats.data?.top_sites)} />
        <RankingList title={t('dashboard.topEventTypes')} data={groupCountToNameCount(eventStats.data?.top_types)} />
        <RankingList title={t('dashboard.topReferrers')} data={referrers} />
        <RankingList title={t('dashboard.topPages')} data={eventStats.data?.top_pages ?? []} />
      </div>

      {/* Distribution pie charts - 3 columns */}
      <div className="grid grid-cols-1 gap-3 lg:grid-cols-3 mb-6">
        <MiniPieChart title={t('dashboard.browsers')} data={browsers} />
        <MiniPieChart title={t('dashboard.operatingSystems')} data={oses} />
        <MiniPieChart title={t('dashboard.languages')} data={languages} />
      </div>

      {/* Bot trend + Hourly - 2 columns */}
      <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
        <TrendChart
          title={t('dashboard.botTrend')}
          data={botTrendData}
          series={[
            { key: 'clickBots', color: '#ef4444', name: t('dashboard.clickBots') },
            { key: 'eventBots', color: '#f59e0b', name: t('dashboard.eventBots') },
          ]}
          height={250}
        />
        <HourlyChart
          title={t('dashboard.hourlyDistribution')}
          clickData={clickStats.data?.hourly ?? []}
          eventData={eventStats.data?.hourly ?? []}
        />
      </div>
    </div>
  )
}
