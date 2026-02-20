import { LogOut } from 'lucide-react'
import { useAuth } from '@/context/AuthContext'
import { Button } from '@/components/ui/Button'

export function Header() {
  const { logout } = useAuth()

  return (
    <header className="flex h-16 items-center justify-between border-b border-border bg-bg-card px-6 shrink-0">
      <h1 className="text-lg font-semibold font-mono text-text">Admin Dashboard</h1>
      <Button variant="ghost" size="sm" onClick={logout}>
        <LogOut className="h-4 w-4" />
        Logout
      </Button>
    </header>
  )
}
