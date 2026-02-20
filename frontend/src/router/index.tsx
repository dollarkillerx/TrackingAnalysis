import { createBrowserRouter } from 'react-router-dom'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import { ProtectedRoute } from './ProtectedRoute'
import { LoginPage } from '@/pages/LoginPage'
import { DashboardPage } from '@/pages/DashboardPage'
import { TrackersPage } from '@/pages/TrackersPage'
import { CampaignsPage } from '@/pages/CampaignsPage'
import { ChannelsPage } from '@/pages/ChannelsPage'
import { TargetsPage } from '@/pages/TargetsPage'
import { SitesPage } from '@/pages/SitesPage'
import { TokenListPage } from '@/pages/TokenListPage'
import { TokenGeneratorPage } from '@/pages/TokenGeneratorPage'
import { NotFoundPage } from '@/pages/NotFoundPage'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <DashboardLayout />
      </ProtectedRoute>
    ),
    children: [
      { index: true, element: <DashboardPage /> },
      { path: 'trackers', element: <TrackersPage /> },
      { path: 'campaigns', element: <CampaignsPage /> },
      { path: 'channels', element: <ChannelsPage /> },
      { path: 'targets', element: <TargetsPage /> },
      { path: 'sites', element: <SitesPage /> },
      { path: 'tokens', element: <TokenListPage /> },
      { path: 'tokens/new', element: <TokenGeneratorPage /> },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])
