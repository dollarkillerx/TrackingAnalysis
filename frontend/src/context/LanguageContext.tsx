import { createContext, useContext, useCallback, useSyncExternalStore } from 'react'
import { useTranslation } from 'react-i18next'
import { supportedLanguages, type SupportedLanguage } from '@/i18n'

interface LanguageContextValue {
  locale: SupportedLanguage
  setLocale: (lang: SupportedLanguage) => void
}

const LanguageContext = createContext<LanguageContextValue | null>(null)

function subscribeToLanguage(callback: () => void) {
  window.addEventListener('languagechange', callback)
  return () => window.removeEventListener('languagechange', callback)
}

export function LanguageProvider({ children }: { children: React.ReactNode }) {
  const { i18n } = useTranslation()

  const locale = useSyncExternalStore(
    subscribeToLanguage,
    () => (supportedLanguages.includes(i18n.language as SupportedLanguage)
      ? (i18n.language as SupportedLanguage)
      : 'en-US'),
  )

  const setLocale = useCallback(
    (lang: SupportedLanguage) => {
      i18n.changeLanguage(lang)
      document.documentElement.lang = lang
    },
    [i18n],
  )

  // Sync html lang on mount
  document.documentElement.lang = locale

  return (
    <LanguageContext.Provider value={{ locale, setLocale }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  const ctx = useContext(LanguageContext)
  if (!ctx) throw new Error('useLanguage must be used within LanguageProvider')
  return ctx
}
