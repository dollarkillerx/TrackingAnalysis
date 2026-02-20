import type enUS from './locales/en-US.json'

export type TranslationKeys = typeof enUS

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'translation'
    resources: {
      translation: TranslationKeys
    }
  }
}
