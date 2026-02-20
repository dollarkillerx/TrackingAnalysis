import { useLanguage } from '@/context/LanguageContext'
import { classNames } from '@/lib/utils'

export function LanguageSwitcher() {
  const { locale, setLocale } = useLanguage()

  return (
    <div className="flex items-center rounded-lg border border-border text-xs">
      <button
        onClick={() => setLocale('en-US')}
        className={classNames(
          'px-2 py-1 rounded-l-lg transition-colors cursor-pointer',
          locale === 'en-US'
            ? 'bg-primary/15 text-primary font-medium'
            : 'text-muted hover:text-text',
        )}
      >
        EN
      </button>
      <button
        onClick={() => setLocale('zh-CN')}
        className={classNames(
          'px-2 py-1 rounded-r-lg transition-colors cursor-pointer',
          locale === 'zh-CN'
            ? 'bg-primary/15 text-primary font-medium'
            : 'text-muted hover:text-text',
        )}
      >
        中文
      </button>
    </div>
  )
}
