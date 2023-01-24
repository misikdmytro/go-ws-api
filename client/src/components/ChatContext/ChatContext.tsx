import React, { useEffect, useState } from 'react'
import useWebSocket from 'react-use-websocket'
import { ReadyState } from 'react-use-websocket/dist/lib/constants'
import { ChatContext as Context, ChatContextProps as ContextProps } from '../../context/chat'
import { Message, WebSocketMessage } from '../../types'

interface ChatContextProps {
  children: JSX.Element | JSX.Element[]
}

export default function ChatContext (props: ChatContextProps): React.ReactElement {
  const { children } = props

  const [context, setContext] = useState<ContextProps>({ messages: [] })
  const { sendMessage, lastMessage, readyState } = useWebSocket('ws://localhost:8080/ws')

  useEffect(() => {
    if (lastMessage !== null) {
      const { type, content }: WebSocketMessage = JSON.parse(lastMessage.data)

      if (type === 'MEMBER_JOIN') {
        const chatMessage: Message = { text: 'New member joined!', sender: content.id, timestamp: 'timestamp' }
        setContext((old) => ({ ...old, messages: [...old.messages, chatMessage] }))
      } else if (type === 'MEMBER_LEAVE') {
        const chatMessage: Message = { text: 'Member leaves the chat', sender: content.id, timestamp: 'timestamp' }
        setContext((old) => ({ ...old, messages: [...old.messages, chatMessage] }))
      }
    }
  }, [lastMessage])

  useEffect(() => {
    if (readyState === ReadyState.OPEN) {
      const msg = {
        type: 'NEW_CLIENT',
        context: {
        }
      }

      sendMessage(JSON.stringify(msg))
    }
  }, [sendMessage, readyState])

  return <Context.Provider value={context}>
        {children}
    </Context.Provider>
}
