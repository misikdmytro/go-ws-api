import React, { useContext } from 'react'
import Message from '../Message'
import { Message as MessageType } from '../../types'
import { ChatContext } from '../../context/chat'

export default function Messages (): React.ReactElement {
  const { messages } = useContext(ChatContext)

  return <>
        {messages.map((msg: MessageType, index: number) => <Message key={index} message={msg} />)}
    </>
}
