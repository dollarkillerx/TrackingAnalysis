import { useEffect, useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useTrackers } from '@/hooks/useTrackers'
import { useCampaigns } from '@/hooks/useCampaigns'
import { useChannels } from '@/hooks/useChannels'
import { useTargets } from '@/hooks/useTargets'
import { useTokens } from '@/hooks/useTokens'
import { useToast } from '@/context/ToastContext'
import { SelectField } from '@/components/ui/SelectField'
import { Button } from '@/components/ui/Button'
import { CopyButton } from '@/components/ui/CopyButton'
import { Key } from 'lucide-react'

export function TokenGeneratorPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { trackers, fetch: fetchTrackers } = useTrackers()
  const { campaigns, fetch: fetchCampaigns } = useCampaigns()
  const { channels, fetch: fetchChannels } = useChannels()
  const { targets, fetch: fetchTargets } = useTargets()
  const { generate: generateToken, generating } = useTokens()
  const { addToast } = useToast()

  const [trackerId, setTrackerId] = useState('')
  const [campaignId, setCampaignId] = useState('')
  const [channelId, setChannelId] = useState('')
  const [targetId, setTargetId] = useState('')
  const [mode, setMode] = useState('302')
  const [generatedShortCode, setGeneratedShortCode] = useState('')
  const [trackingUrl, setTrackingUrl] = useState('')

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
      })
      setGeneratedShortCode(result.short_code)
      setTrackingUrl(result.tracking_url)
      addToast(t('tokens.tokenGenerated'), 'success')
    } catch (err) {
      addToast(err instanceof Error ? err.message : t('tokens.failedToGenerate'), 'error')
    }
  }

  return (
    <div>
      <h2 className="mb-6 text-2xl font-bold font-mono text-text">{t('tokens.title')}</h2>

      <div className="max-w-xl rounded-xl border border-border bg-bg-card p-6">
        <form onSubmit={handleGenerate} className="space-y-4">
          <SelectField
            label={t('common.tracker')}
            value={trackerId}
            onChange={(e) => setTrackerId(e.target.value)}
            options={trackers.map((t) => ({ value: t.id, label: `${t.name} (${t.type})` }))}
            placeholder={t('common.selectTracker')}
            required
          />
          <SelectField
            label={t('common.campaign')}
            value={campaignId}
            onChange={(e) => setCampaignId(e.target.value)}
            options={filteredCampaigns.map((c) => ({ value: c.id, label: c.name }))}
            placeholder={t('common.selectCampaign')}
          />
          <SelectField
            label={t('common.channel')}
            value={channelId}
            onChange={(e) => setChannelId(e.target.value)}
            options={filteredChannels.map((c) => ({ value: c.id, label: c.name }))}
            placeholder={t('common.selectChannel')}
          />
          <SelectField
            label={t('common.target')}
            value={targetId}
            onChange={(e) => setTargetId(e.target.value)}
            options={filteredTargets.map((t) => ({ value: t.id, label: t.url }))}
            placeholder={t('common.selectTarget')}
            required
          />
          <SelectField
            label={t('tokens.mode')}
            value={mode}
            onChange={(e) => setMode(e.target.value)}
            options={[
              { value: '302', label: t('common.redirect302') },
              { value: 'js', label: t('common.javascript') },
            ]}
          />

          <Button type="submit" loading={generating} className="w-full">
            <Key className="h-4 w-4" /> {t('tokens.generateToken')}
          </Button>
        </form>

        {generatedShortCode && (
          <div className="mt-6 space-y-3">
            <div>
              <label className="block text-sm font-medium text-muted mb-1">{t('tokens.shortCode')}</label>
              <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-deep px-4 py-3">
                <code className="flex-1 font-mono text-xs text-primary break-all">{generatedShortCode}</code>
                <CopyButton text={generatedShortCode} />
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-muted mb-1">{t('tokens.trackingUrl')}</label>
              <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-deep px-4 py-3">
                <code className="flex-1 font-mono text-xs text-secondary break-all">{trackingUrl}</code>
                <CopyButton text={trackingUrl} />
              </div>
            </div>
            <Button variant="ghost" className="w-full" onClick={() => navigate('/tokens')}>
              {t('tokens.viewAllTokens')}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
