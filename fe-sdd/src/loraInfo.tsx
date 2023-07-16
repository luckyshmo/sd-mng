import { useEffect, useState } from 'react';

interface LoraInfo {
  name: string
  token: string
}

const itemStyles = {
  display: 'flex',
  alignItems: 'center',
};

const textStyles = {
  marginRight: '10px',
};

const Item = ({ item }: {item: LoraInfo}) => {
  const [buttonText, setButtonText] = useState('COPY');

  const copyToClipboard = (token: string) => {
    navigator.clipboard.writeText(token)
      .then(() => {
        console.log('Copied to clipboard:', token);
        setButtonText('✔✔✔');

        setTimeout(() => {
          setButtonText('Copy');
        }, 3000); // Reset button text after 3 seconds
      })
      .catch((error) => {
        console.error('Failed to copy to clipboard:', error);
      });
  };

  return (
    <div style={itemStyles}>
      <p style={textStyles}>
        {item.name}: {item.token}
      </p>
      <button onClick={() => copyToClipboard(item.token)}>{buttonText}</button>
    </div>
  );
};

const LoraInfoComponent = () => {
  const [items, setItems] = useState<LoraInfo[]>([]);
  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
      .then(() => {
        console.log('Text copied to clipboard:', text);
      })
      .catch((error) => {
        console.error('Failed to copy text to clipboard:', error);
      });
  }

  useEffect(() => {
    const address = 'http://192.168.1.10:8080/info/lora';
    fetch(address, {method: 'GET'})
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.text();
      })
      .then((responseData) => {
        let info: LoraInfo[] = JSON.parse(responseData)
        setItems(info)
        console.log(responseData);
      })
      .catch((error) => {
        console.error(`Error: ${error.message}`);
      });   
  }, []);

  if (items.length > 0) {
    return <div>
      <h2>LoraInfo</h2>
      <div>
      {items.map((item, index) => {
        return <Item key={index} item={item} />;
      })}
      </div>
    </div>
  }

  return <div></div>
};

export default LoraInfoComponent;
