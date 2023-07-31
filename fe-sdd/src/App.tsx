import React from 'react'
import WebSocketConn from './api/ws'
import { useEffect } from 'react'
import NavBar from './components/nav'
import { BrowserRouter } from 'react-router-dom'

const App: React.FC = () => {
  useEffect(() => {
    const webSocketConnection = new WebSocketConn()
    webSocketConnection.connect()

    return () => {
      // Close the WebSocket connection on component unmount
      if (webSocketConnection.socket) {
        webSocketConnection.socket.close()
      }
    }
  }, [])

  return (
    <BrowserRouter>
      <NavBar />
    </BrowserRouter>
  )
}

export default App
