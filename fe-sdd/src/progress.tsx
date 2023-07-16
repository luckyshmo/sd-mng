import { useEffect, useState } from 'react';

interface item {
  id: string
  name: string
  percentage: number
}

const ProgressComponent = () => {
  const [items, setItems] = useState<item[]>([]);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  const upd = function updateById(id: string, updatedItem: item): void {
    const itemIndex = items.findIndex((item) => item.id === id);
    if (itemIndex !== -1) {
      items[itemIndex] = updatedItem
    } else {
      items.push(updatedItem)
    }
    let newItems = [...items]
    setItems(newItems)
  }

  useEffect(() => {
    const ws = new WebSocket('ws://192.168.1.10:8080/progress');
    console.log(ws.url)
    setSocket(ws);

    ws.onopen = () => {
      console.log('WebSocket connection established.');
    };

    ws.onmessage = (event) => {
      console.log("event: ", event.data)
      let msg: item =  JSON.parse(event.data)
      upd(msg.id, msg)
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed.');
    };

    return () => {
      // Clean up the WebSocket connection on component unmount
      if (socket) {
        socket.close();
      }
    };
  }, []);

  if (items.length > 0) {
    return <div>
      <h2>Progress</h2>
      <div>
      {items.map((item, index) => {
        return <div key={index}>
          <p>{item.name}: {item.percentage}%</p>
        </div>;
      })}
      </div>
    </div>
  }

  return <div></div>
};

export default ProgressComponent;
