import { createContext } from 'react'
import { Message } from '../types'

export interface ChatContextProps {
  messages: Message[]
}

export const ChatContext = createContext<ChatContextProps>({ messages: [] })
