export interface Tracker {
  id: string
  type: 'ad' | 'web'
  name: string
  status: string
  created_at: string
  updated_at: string
}

export interface Campaign {
  id: string
  tracker_id: string
  name: string
  status: string
  created_at: string
  updated_at: string
}

export interface Channel {
  id: string
  tracker_id: string
  campaign_id: string
  name: string
  source: string
  medium: string
  tags: Record<string, string>
  created_at: string
  updated_at: string
}

export interface Target {
  id: string
  tracker_id: string
  url: string
  created_at: string
}

export interface Site {
  id: string
  name: string
  domain: string
  site_key: string
  status: string
  created_at: string
  updated_at: string
}

export interface Token {
  id: string
  short_code: string
  tracker_id: string
  campaign_id: string
  channel_id: string
  target_id: string
  mode: string
  created_at: string
  tracking_url: string
}

export interface DailyCount {
  date: string
  count: number
}

export interface GroupCount {
  group_id: string
  name: string
  count: number
}

export interface HourlyCount {
  hour: number
  count: number
}

export interface NameCount {
  name: string
  count: number
}

export interface ClickStatsResponse {
  summary: { total: number; unique_visitors: number; bots: number; bot_rate: number }
  daily: DailyCount[]
  top_trackers: GroupCount[]
  top_channels: GroupCount[]
  top_campaigns: GroupCount[]
  top_referrers: NameCount[]
  browsers: NameCount[]
  oses: NameCount[]
  languages: NameCount[]
  bot_daily: DailyCount[]
  hourly: HourlyCount[]
}

export interface EventStatsResponse {
  summary: { total: number; unique_visitors: number; unique_sessions: number; bots: number; bot_rate: number }
  daily: DailyCount[]
  top_sites: GroupCount[]
  top_types: GroupCount[]
  top_referrers: NameCount[]
  top_pages: NameCount[]
  browsers: NameCount[]
  oses: NameCount[]
  languages: NameCount[]
  bot_daily: DailyCount[]
  hourly: HourlyCount[]
}
