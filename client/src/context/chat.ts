import { createContext } from 'react'
import { Message } from '../types'

export interface ChatContextProps {
  messages: Message[]
  sendMessage: (text: string) => void
  id: string
}

export const ChatContext = createContext<ChatContextProps>({ messages: [], sendMessage: () => {}, id: '' })
