import { userHttp } from './http'
import type {
  LoginReq,
  LoginResp,
  RegisterReq,
  UserInfoResp,
} from '@/types/api'

export function login(body: LoginReq): Promise<LoginResp> {
  return userHttp.post('/login', body)
}

export function register(body: RegisterReq): Promise<LoginResp> {
  return userHttp.post('/register', body)
}

export function getUserInfo(userId: string): Promise<UserInfoResp> {
  return userHttp.get('/user', { params: { user_id: userId } })
}
