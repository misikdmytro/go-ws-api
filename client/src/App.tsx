import React from 'react'
import { useWebSocket } from './hooks/ws'

function App (): React.ReactElement<any, any> {
  useWebSocket('ws://localhost:8080/ws')
  return <div>Hello, world!</div>
}

export default App
