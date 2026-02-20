export function formatDate(dateStr: string, locale: string = 'en-US'): string {
  return new Date(dateStr).toLocaleDateString(locale, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function truncate(str: string, len = 40): string {
  if (str.length <= len) return str
  return str.slice(0, len) + '...'
}

export async function copyToClipboard(text: string): Promise<void> {
  await navigator.clipboard.writeText(text)
}

export function classNames(...classes: (string | false | null | undefined)[]): string {
  return classes.filter(Boolean).join(' ')
}

export function mergeNameCounts(...arrays: ({ name: string; count: number }[] | undefined)[]): { name: string; count: number }[] {
  const map = new Map<string, number>()
  for (const arr of arrays) {
    if (!arr) continue
    for (const item of arr) {
      map.set(item.name, (map.get(item.name) ?? 0) + item.count)
    }
  }
  return Array.from(map.entries())
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
}
