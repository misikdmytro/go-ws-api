import React from 'react'
import Box from '@mui/material/Box'
import { Message as MessageType } from '../../types'

interface MesssageProps {
  message: MessageType
}

export default function Message (props: MesssageProps): React.ReactElement {
  const { message: { sender, text } } = props
  return <Box>
    <b>{sender}</b>: {text}
  </Box>
}
