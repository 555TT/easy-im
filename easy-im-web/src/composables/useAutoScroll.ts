import { ref, type Ref, nextTick } from 'vue'

export function useAutoScroll(container: Ref<HTMLElement | null>) {
  const nearBottom = ref(true)
  const threshold = 40

  function onScroll(): void {
    const el = container.value
    if (!el) return
    nearBottom.value = el.scrollHeight - el.scrollTop - el.clientHeight < threshold
  }

  async function scrollToBottom(force = false): Promise<void> {
    await nextTick()
    const el = container.value
    if (!el) return
    if (force || nearBottom.value) {
      el.scrollTop = el.scrollHeight
    }
  }

  return { nearBottom, onScroll, scrollToBottom }
}
