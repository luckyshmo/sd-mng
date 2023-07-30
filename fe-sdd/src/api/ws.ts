import { DownloadModel, downloadStore } from "../store/store";

const WS_URL =  import.meta.env.VITE_WS_URL!

class WebSocketConn {
  socket: WebSocket | null = null;
  reconnectInterval: number = 5000; //! conf
  reconnectAttempt: number = 0;
  maxReconnectAttempts: number = 5;

  connect() {
    this.socket = new WebSocket(WS_URL);

    this.socket.onopen = () => {
      console.log('WebSocket connection established.');
    };

    this.socket.onmessage = (event) => {
      console.log("event: ", event.data)
      let msg: DownloadModel =  JSON.parse(event.data)
      downloadStore.updDownload(msg.id, msg)
    };

    this.socket.onclose = () => {
      console.log('WebSocket connection closed.');
      this.reconnect();
    };
  }

  private reconnect() {
    if (this.reconnectAttempt < this.maxReconnectAttempts) {
      console.log(`Trying to reconnect... Attempt ${this.reconnectAttempt + 1}`);
      setTimeout(() => {
        this.reconnectAttempt++;
        this.connect();
      }, this.reconnectInterval);
    } else {
      console.log('Max reconnect attempts reached. Cannot establish WebSocket connection.');
    }
  }
}

export default WebSocketConn;