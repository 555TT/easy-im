import axios, { type AxiosInstance, type AxiosError } from 'axios'

export class ApiError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message)
    this.code = code
    this.name = 'ApiError'
  }
}

let getToken: () => string | null = () => null
let onUnauthorized: () => void = () => {}

export function configureHttp(opts: {
  getToken: () => string | null
  onUnauthorized: () => void
}): void {
  getToken = opts.getToken
  onUnauthorized = opts.onUnauthorized
}

function createInstance(baseURL: string): AxiosInstance {
  const i = axios.create({ baseURL, timeout: 10_000 })

  i.interceptors.request.use((cfg) => {
    const t = getToken()
    if (t) cfg.headers.Authorization = `Bearer ${t}`
    return cfg
  })

  i.interceptors.response.use(
    (resp) => {
      const body = resp.data
      if (body && typeof body === 'object' && 'code' in body) {
        if (body.code === 200) return body.data
        if (body.code === 401) onUnauthorized()
        throw new ApiError(body.code, body.msg ?? 'request failed')
      }
      return body
    },
    (err: AxiosError) => {
      if (err.response?.status === 401) onUnauthorized()
      const body = err.response?.data as
        | { code?: number; msg?: string }
        | undefined
      throw new ApiError(
        body?.code ?? err.response?.status ?? 0,
        body?.msg ?? err.message,
      )
    },
  )

  return i
}

export const userHttp = createInstance('/api/user')
export const socialHttp = createInstance('/api/social')
export const imHttp = createInstance('/api/im')
