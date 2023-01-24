import React, { useEffect, useState } from 'react'
import useWebSocket from 'react-use-websocket'
import { ChatContext as Context, ChatContextProps as ContextProps } from '../../context/chat'
import { WebSocketMessage } from '../../types'

interface ChatContextProps {
  children: JSX.Element | JSX.Element[]
}

export default function ChatContext (props: ChatContextProps): React.ReactElement {
  const { children } = props

  const { sendMessage: sendWsMessage, lastMessage } = useWebSocket('ws://localhost:8080/ws')
  const sendMessage = (text: string): void => {
    const msg: WebSocketMessage = {
      type: 'MESSAGE',
      content: {
        message: text
      }
    }
    sendWsMessage(JSON.stringify(msg))
    setContext((old) => ({ ...old, messages: [...old.messages, { text, sender: 'You' }] }))
  }
  const [context, setContext] = useState<ContextProps>({ messages: [], sendMessage, id: '' })

  useEffect(() => {
    if (lastMessage !== null) {
      const { type, content }: WebSocketMessage = JSON.parse(lastMessage.data)

      if (type === 'ID_ASSIGNED') {
        setContext((old) => ({ ...old, id: content.id }))
      } else if (type === 'MEMBER_JOIN') {
        setContext((old) => ({ ...old, messages: [...old.messages, { text: 'New member joined!', sender: content.id }] }))
      } else if (type === 'MEMBER_LEAVE') {
        setContext((old) => ({ ...old, messages: [...old.messages, { text: 'Member leaves the chat', sender: content.id }] }))
      } else if (type === 'MESSAGE') {
        setContext((old) => ({ ...old, messages: [...old.messages, { text: content.message, sender: content.id }] }))
      }
    }
  }, [lastMessage])

  return <Context.Provider value={context}>
        {children}
    </Context.Provider>
}
