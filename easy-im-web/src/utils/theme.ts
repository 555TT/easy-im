export type ThemeMode = 'light' | 'dark'

const THEME_KEY = 'easy-im.theme'

export function getStoredTheme(): ThemeMode {
  const theme = localStorage.getItem(THEME_KEY)
  return theme === 'dark' ? 'dark' : 'light'
}

export function applyTheme(theme: ThemeMode): void {
  document.documentElement.setAttribute('data-theme', theme)
}

export function setStoredTheme(theme: ThemeMode): void {
  localStorage.setItem(THEME_KEY, theme)
  applyTheme(theme)
}

export function applyStoredTheme(): ThemeMode {
  const theme = getStoredTheme()
  applyTheme(theme)
  return theme
}
