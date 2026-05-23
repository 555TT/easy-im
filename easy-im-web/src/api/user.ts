import { userHttp } from './http'
import type {
  LoginReq,
  LoginResp,
  RegisterReq,
  UpdateProfileReq,
  UpdateProfileResp,
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

export function updateProfile(body: UpdateProfileReq): Promise<UpdateProfileResp> {
  return userHttp.put('/profile', body)
}
