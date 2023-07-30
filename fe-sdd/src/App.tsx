import React from 'react';
import UrlInput from './UrlInput';
import ProgressComponent from './progress';
import LoraInfoComponent from './loraInfo';
import WebSocketConn from './api/ws';
import { useEffect } from 'react';

const App: React.FC = () => {
  useEffect(() => {
    const webSocketConnection = new WebSocketConn();
    webSocketConnection.connect();

    return () => {
      // Close the WebSocket connection on component unmount
      if (webSocketConnection.socket) {
        webSocketConnection.socket.close();
      }
    };
  }, []);

  return (
    <div>
      <h1>Request Sender</h1>
      <UrlInput />
      <LoraInfoComponent />
      <ProgressComponent />
    </div>
  );
};

export default App;
