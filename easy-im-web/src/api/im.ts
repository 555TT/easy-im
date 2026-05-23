import { imHttp } from './http'
import type {
  ChatLogListResp,
  ChatLogReq,
  GetConversationsResp,
  PutConversationsReq,
  SetUpUserConversationReq,
} from '@/types/api'

export function listConversations(): Promise<GetConversationsResp> {
  return imHttp.get('/conversation')
}

export function putConversations(body: PutConversationsReq): Promise<unknown> {
  return imHttp.put('/conversation', body)
}

export function setUpConversation(body: SetUpUserConversationReq): Promise<unknown> {
  return imHttp.post('/setup/conversation', body)
}

export function listChatLog(params: ChatLogReq): Promise<ChatLogListResp> {
  return imHttp.get('/chatlog', { params })
}

export function listChatLogReadRecords(
  msgId: string,
): Promise<{ reads: string[]; unReads: string[] }> {
  return imHttp.get('/chatlog/readRecords', { params: { msgId } })
}
