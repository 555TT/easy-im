let counter = 0

export function clientMsgId(): string {
  counter = (counter + 1) % 1_000_000
  return `c-${Date.now()}-${counter}`
}
