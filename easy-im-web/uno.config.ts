import { defineConfig, presetIcons, presetUno } from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetIcons({ scale: 1.2, cdn: 'https://esm.sh/' }),
  ],
  theme: {
    colors: {
      brand: {
        400: '#4096ff',
        500: '#1677ff',
        600: '#0958d9',
      },
      rail: {
        bg: '#1f2733',
        active: '#1677ff',
      },
      side: {
        bg: '#f5f7fa',
        hover: '#eef0f3',
        active: '#e6f4ff',
      },
      bubble: {
        mine: '#1677ff',
        theirs: '#f4f5f7',
      },
      ink: {
        primary: '#1f2329',
        secondary: '#646a73',
      },
      ok: '#52c41a',
      danger: '#ff4d4f',
      divider: '#eaecef',
    },
    fontFamily: {
      sans: '-apple-system, BlinkMacSystemFont, "PingFang SC", "Microsoft YaHei", "Helvetica Neue", Arial, sans-serif',
    },
  },
  shortcuts: {
    'flex-center': 'flex items-center justify-center',
    'flex-between': 'flex items-center justify-between',
  },
})
