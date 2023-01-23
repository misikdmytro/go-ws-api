import { useEffect, useState } from 'react'

export const useWebSocket = (addr: string): WebSocket => {
  const [ws] = useState(() => new WebSocket(addr, 'echo-protocol'))

  useEffect(() => {
    ws.addEventListener('open', () => {
      const msg = {
        type: 'type',
        content: {}
      }

      ws.send(JSON.stringify(msg))
    })
    ws.addEventListener('message', () => { console.log('message') })

    return () => { ws.close() }
  }, [])

  return ws
}
