import React, { useContext, useState } from 'react'
import { Box, Button, SxProps, TextField, Theme } from '@mui/material'
import SendIcon from '@mui/icons-material/Send'
import { ChatContext } from '../../context/chat'

interface TextInputProps {
  sx?: SxProps<Theme>
}

export default function TextInput (props: TextInputProps): React.ReactElement {
  const { sendMessage } = useContext(ChatContext)
  const [value, setValue] = useState('')

  const onClick = (): void => {
    sendMessage(value)
    setValue('')
  }

  return (
        <Box sx={props.sx}>
            <form style={{
              display: 'flex',
              justifyContent: 'center'
            }} noValidate autoComplete="off">
                <TextField
                    label="Message"
                    value={value}
                    onChange={(ev) => { setValue(ev.target.value) }}
                />
                <Button variant="contained" color="primary" onClick={onClick}>
                    <SendIcon />
                </Button>
            </form>
        </Box>
  )
}
