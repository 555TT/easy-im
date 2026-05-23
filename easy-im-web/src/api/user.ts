import { userHttp } from './http'
import type {
  LoginReq,
  LoginResp,
  RegisterReq,
  UpdateProfileReq,
  UpdateProfileResp,
  UpdatePasswordReq,
  UpdatePasswordResp,
  UserInfoResp,
} from '@/types/api'

export function login(body: LoginReq): Promise<LoginResp> {
  return userHttp.post('/login', body)
}

export function register(body: RegisterReq): Promise<LoginResp> {
  return userHttp.post('/register', body)
}

export function getUserInfo(): Promise<UserInfoResp> {
  return userHttp.get('/user')
}

export function updateProfile(body: UpdateProfileReq): Promise<UpdateProfileResp> {
  return userHttp.put('/profile', body)
}

export function updatePassword(body: UpdatePasswordReq): Promise<UpdatePasswordResp> {
  return userHttp.put('/password', body)
}
