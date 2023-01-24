import { Box } from '@mui/material'
import Paper from '@mui/material/Paper'
import React from 'react'
import ChatContext from '../ChatContext'
import Messages from '../Messages/Messages'
import TextInput from '../TextInput'

export default function Chat (): React.ReactElement {
  return <ChatContext>
    <Box sx={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      margin: 4
    }}>
        <Paper sx={{
          width: '80vw',
          height: '80vh',
          maxWidth: '500px',
          maxHeight: '700px',
          display: 'flex',
          alignItems: 'center',
          flexDirection: 'column',
          position: 'relative'
        }}>
            <Paper sx={{
              width: '100%',
              overflowY: 'scroll',
              height: 'calc( 100% - 80px )'
            }}>
                <Messages />
            </Paper>
        <TextInput sx={{
          width: '100%',
          margin: 2
        }} />
        </Paper>
    </Box>
  </ChatContext>
}
