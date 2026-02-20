import { useEffect } from 'react'
import { Crosshair, Megaphone, Globe, Share2 } from 'lucide-react'
import { useTrackers } from '@/hooks/useTrackers'
import { useCampaigns } from '@/hooks/useCampaigns'
import { useSites } from '@/hooks/useSites'
import { useChannels } from '@/hooks/useChannels'
import type { LucideIcon } from 'lucide-react'

function StatCard({ icon: Icon, label, count, loading }: { icon: LucideIcon; label: string; count: number; loading: boolean }) {
  return (
    <div className="rounded-xl border border-border bg-bg-card p-6">
      <div className="flex items-center gap-4">
        <div className="rounded-lg bg-primary/10 p-3">
          <Icon className="h-6 w-6 text-primary" />
        </div>
        <div>
          <p className="text-2xl font-bold font-mono text-text">
            {loading ? '...' : count}
          </p>
          <p className="text-sm text-muted">{label}</p>
        </div>
      </div>
    </div>
  )
}

export function DashboardPage() {
  const { trackers, loading: tl, fetch: fetchTrackers } = useTrackers()
  const { campaigns, loading: cl, fetch: fetchCampaigns } = useCampaigns()
  const { sites, loading: sl, fetch: fetchSites } = useSites()
  const { channels, loading: chl, fetch: fetchChannels } = useChannels()

  useEffect(() => {
    fetchTrackers()
    fetchCampaigns()
    fetchSites()
    fetchChannels()
  }, [fetchTrackers, fetchCampaigns, fetchSites, fetchChannels])

  return (
    <div>
      <h2 className="mb-6 text-2xl font-bold font-mono text-text">Dashboard</h2>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard icon={Crosshair} label="Trackers" count={trackers.length} loading={tl} />
        <StatCard icon={Megaphone} label="Campaigns" count={campaigns.length} loading={cl} />
        <StatCard icon={Globe} label="Sites" count={sites.length} loading={sl} />
        <StatCard icon={Share2} label="Channels" count={channels.length} loading={chl} />
      </div>
    </div>
  )
}
