import React from 'react'
import { Box, Button, SxProps, TextField, Theme } from '@mui/material'
import SendIcon from '@mui/icons-material/Send'

interface TextInputProps {
  sx?: SxProps<Theme>
}

export default function TextInput (props: TextInputProps): React.ReactElement {
  return (
        <Box sx={props.sx}>
            <form style={{
              display: 'flex',
              justifyContent: 'center'
            }} noValidate autoComplete="off">
                <TextField
                    label="Message"
                />
                <Button variant="contained" color="primary">
                    <SendIcon />
                </Button>
            </form>
        </Box>
  )
}
