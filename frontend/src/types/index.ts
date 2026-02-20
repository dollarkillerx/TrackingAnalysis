export interface Tracker {
  id: string
  type: 'ad' | 'web'
  name: string
  mode: string
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
