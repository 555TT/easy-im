export function formatHM(ts: number): string {
  if (!ts) return ''
  const d = new Date(ts)
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${hh}:${mm}`
}

export function formatRelative(ts: number, now: number = Date.now()): string {
  if (!ts) return ''
  const d = new Date(ts)
  const n = new Date(now)
  const sameDay =
    d.getFullYear() === n.getFullYear() &&
    d.getMonth() === n.getMonth() &&
    d.getDate() === n.getDate()
  if (sameDay) return formatHM(ts)
  const diffDays = Math.floor((now - ts) / 86_400_000)
  if (diffDays < 7) {
    return ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][d.getDay()]
  }
  const mo = String(d.getMonth() + 1).padStart(2, '0')
  const da = String(d.getDate()).padStart(2, '0')
  return `${mo}-${da}`
}
