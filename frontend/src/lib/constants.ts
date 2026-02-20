export const RPC = {
  LOGIN: 'admin.login',
  TRACKER_CREATE: 'admin.tracker.create',
  TRACKER_LIST: 'admin.tracker.list',
  TRACKER_UPDATE: 'admin.tracker.update',
  TRACKER_DELETE: 'admin.tracker.delete',
  CAMPAIGN_CREATE: 'admin.campaign.create',
  CAMPAIGN_LIST: 'admin.campaign.list',
  CHANNEL_CREATE: 'admin.channel.create',
  CHANNEL_BATCH_IMPORT: 'admin.channel.batchImport',
  CHANNEL_LIST: 'admin.channel.list',
  TARGET_CREATE: 'admin.target.create',
  TARGET_LIST: 'admin.target.list',
  SITE_CREATE: 'admin.site.create',
  SITE_LIST: 'admin.site.list',
  TOKEN_GENERATE: 'admin.token.generate',
  TOKEN_LIST: 'admin.token.list',
  TOKEN_DELETE: 'admin.token.delete',
  STATS_CLICKS: 'admin.stats.clicks',
  STATS_EVENTS: 'admin.stats.events',
} as const

export const TOKEN_KEY = 'admin_token'
