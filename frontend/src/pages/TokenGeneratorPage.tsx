import { useEffect, useState, type FormEvent } from 'react'
import { useTrackers } from '@/hooks/useTrackers'
import { useCampaigns } from '@/hooks/useCampaigns'
import { useChannels } from '@/hooks/useChannels'
import { useTargets } from '@/hooks/useTargets'
import { useRpcCall } from '@/hooks/useRpc'
import { useToast } from '@/context/ToastContext'
import { SelectField } from '@/components/ui/SelectField'
import { Button } from '@/components/ui/Button'
import { CopyButton } from '@/components/ui/CopyButton'
import { RPC } from '@/lib/constants'
import { Key } from 'lucide-react'

const expiryOptions = [
  { value: '3600', label: '1 hour' },
  { value: '86400', label: '24 hours' },
  { value: '604800', label: '7 days' },
  { value: '2592000', label: '30 days' },
  { value: '0', label: 'No expiry' },
]

export function TokenGeneratorPage() {
  const { trackers, fetch: fetchTrackers } = useTrackers()
  const { campaigns, fetch: fetchCampaigns } = useCampaigns()
  const { channels, fetch: fetchChannels } = useChannels()
  const { targets, fetch: fetchTargets } = useTargets()
  const { execute: generateToken, loading } = useRpcCall<{ token: string }>(RPC.TOKEN_GENERATE)
  const { addToast } = useToast()

  const [trackerId, setTrackerId] = useState('')
  const [campaignId, setCampaignId] = useState('')
  const [channelId, setChannelId] = useState('')
  const [targetId, setTargetId] = useState('')
  const [mode, setMode] = useState('302')
  const [expSeconds, setExpSeconds] = useState('86400')
  const [generatedToken, setGeneratedToken] = useState('')

  useEffect(() => { fetchTrackers() }, [fetchTrackers])

  useEffect(() => {
    if (trackerId) {
      fetchCampaigns({ tracker_id: trackerId })
      fetchTargets({ tracker_id: trackerId })
      setCampaignId('')
      setChannelId('')
      setTargetId('')
    }
  }, [trackerId, fetchCampaigns, fetchTargets])

  useEffect(() => {
    if (campaignId) {
      fetchChannels({ tracker_id: trackerId, campaign_id: campaignId })
      setChannelId('')
    }
  }, [campaignId, trackerId, fetchChannels])

  const filteredCampaigns = campaigns.filter((c) => c.tracker_id === trackerId)
  const filteredChannels = channels.filter((c) => c.campaign_id === campaignId)
  const filteredTargets = targets.filter((t) => t.tracker_id === trackerId)

  const handleGenerate = async (e: FormEvent) => {
    e.preventDefault()
    try {
      const result = await generateToken({
        tracker_id: trackerId,
        campaign_id: campaignId,
        channel_id: channelId,
        target_id: targetId,
        mode,
        exp_seconds: Number(expSeconds),
      })
      setGeneratedToken(result.token)
      addToast('Token generated', 'success')
    } catch (err) {
      addToast(err instanceof Error ? err.message : 'Failed to generate token', 'error')
    }
  }

  const trackingUrl = generatedToken
    ? `${window.location.origin}/${mode === 'js' ? 't' : 'r'}/${generatedToken}`
    : ''

  return (
    <div>
      <h2 className="mb-6 text-2xl font-bold font-mono text-text">Token Generator</h2>

      <div className="max-w-xl rounded-xl border border-border bg-bg-card p-6">
        <form onSubmit={handleGenerate} className="space-y-4">
          <SelectField
            label="Tracker"
            value={trackerId}
            onChange={(e) => setTrackerId(e.target.value)}
            options={trackers.map((t) => ({ value: t.id, label: `${t.name} (${t.type})` }))}
            placeholder="Select tracker"
            required
          />
          <SelectField
            label="Campaign"
            value={campaignId}
            onChange={(e) => setCampaignId(e.target.value)}
            options={filteredCampaigns.map((c) => ({ value: c.id, label: c.name }))}
            placeholder="Select campaign"
          />
          <SelectField
            label="Channel"
            value={channelId}
            onChange={(e) => setChannelId(e.target.value)}
            options={filteredChannels.map((c) => ({ value: c.id, label: c.name }))}
            placeholder="Select channel"
          />
          <SelectField
            label="Target"
            value={targetId}
            onChange={(e) => setTargetId(e.target.value)}
            options={filteredTargets.map((t) => ({ value: t.id, label: t.url }))}
            placeholder="Select target"
            required
          />
          <SelectField
            label="Mode"
            value={mode}
            onChange={(e) => setMode(e.target.value)}
            options={[
              { value: '302', label: '302 Redirect' },
              { value: 'js', label: 'JavaScript' },
            ]}
          />
          <SelectField
            label="Expiry"
            value={expSeconds}
            onChange={(e) => setExpSeconds(e.target.value)}
            options={expiryOptions}
          />

          <Button type="submit" loading={loading} className="w-full">
            <Key className="h-4 w-4" /> Generate Token
          </Button>
        </form>

        {generatedToken && (
          <div className="mt-6 space-y-3">
            <div>
              <label className="block text-sm font-medium text-muted mb-1">Token</label>
              <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-deep px-4 py-3">
                <code className="flex-1 font-mono text-xs text-primary break-all">{generatedToken}</code>
                <CopyButton text={generatedToken} />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-muted mb-1">Tracking URL</label>
              <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-deep px-4 py-3">
                <code className="flex-1 font-mono text-xs text-secondary break-all">{trackingUrl}</code>
                <CopyButton text={trackingUrl} />
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
