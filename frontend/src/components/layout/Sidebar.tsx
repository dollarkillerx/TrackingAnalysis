import { useState, useEffect } from 'react'
import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Crosshair,
  Megaphone,
  Share2,
  Target,
  Globe,
  Key,
  PanelLeftClose,
  PanelLeft,
} from 'lucide-react'
import { classNames } from '@/lib/utils'

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/trackers', icon: Crosshair, label: 'Trackers' },
  { to: '/campaigns', icon: Megaphone, label: 'Campaigns' },
  { to: '/channels', icon: Share2, label: 'Channels' },
  { to: '/targets', icon: Target, label: 'Targets' },
  { to: '/sites', icon: Globe, label: 'Sites' },
  { to: '/tokens', icon: Key, label: 'Token Generator' },
]

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(() => window.innerWidth < 768)

  useEffect(() => {
    const onResize = () => {
      if (window.innerWidth < 768) setCollapsed(true)
    }
    window.addEventListener('resize', onResize)
    return () => window.removeEventListener('resize', onResize)
  }, [])

  return (
    <aside
      className={classNames(
        'flex flex-col border-r border-border bg-bg-card transition-all duration-200 shrink-0',
        collapsed ? 'w-16' : 'w-60',
      )}
    >
      <div className="flex h-16 items-center justify-between border-b border-border px-3">
        {!collapsed && (
          <span className="text-sm font-semibold text-text truncate">TrackingAnalysis</span>
        )}
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="rounded-lg p-2 text-muted hover:bg-bg-elevated hover:text-text transition-colors cursor-pointer"
        >
          {collapsed ? <PanelLeft className="h-5 w-5" /> : <PanelLeftClose className="h-5 w-5" />}
        </button>
      </div>

      <nav className="flex-1 overflow-y-auto py-3 px-2 space-y-1">
        {navItems.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              classNames(
                'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors duration-150',
                isActive
                  ? 'bg-primary/10 text-primary border-l-2 border-primary'
                  : 'text-muted hover:bg-bg-elevated hover:text-text',
                collapsed && 'justify-center px-0',
              )
            }
            title={collapsed ? label : undefined}
          >
            <Icon className="h-5 w-5 shrink-0" />
            {!collapsed && <span className="truncate">{label}</span>}
          </NavLink>
        ))}
      </nav>
    </aside>
  )
}
