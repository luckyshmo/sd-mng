import React, { useState, ChangeEvent, FormEvent } from 'react';

type Option = {
  value: string;
  label: string;
};

const options: Option[] = [
  { value: 'models/Lora', label: 'Lora' },
  { value: 'models/Stable-Diffusion', label: 'Model' },
];


const UrlInput: React.FC = () => {
  const [url, setUrl] = useState('');
  const [selectedOption, setSelectedOption] = useState(options[0].value);


  const handleFolderChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedOption(event.target.value);
  };

  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    setUrl(event.target.value);
  };

  function isValidUrl(url: string): boolean {
    // Regular expression to validate URLs
    const urlRegex = new RegExp(
      "^(https?:\\/\\/)?" + // protocol
      "((([a-zA-Z\\d]([a-zA-Z\\d-]{0,61}[a-zA-Z\\d])?)\\.)+[a-zA-Z]{2,}|" + // domain name
      "((\\d{1,3}\\.){3}\\d{1,3}))" + // OR ip (v4) address
      "(\\:\\d+)?(\\/[-a-zA-Z\\d%@_.~+&:]*)*" + // port and path
      "(\\?[;&a-zA-Z\\d%@_.,~+&:=-]*)?" + // query string
      "(\\#[-a-zA-Z\\d_]*)?$", // fragment locator
      "i"
    );
  
    return urlRegex.test(url);
  }

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    
    if (!isValidUrl(url)) {
      alert("Invalid URL")
      return
    }

    
    const address = 'http://192.168.1.10:8080/';
    const urlParams = new URLSearchParams({
      url: url,
      folder: selectedOption,
    });
    fetch(address + '?' + urlParams.toString(), {method: 'POST'})
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.text();
      })
      .then((responseData) => {
        console.log(responseData);
      })
      .catch((error) => {
        console.error(`Error: ${error.message}`);
      });
  };

  return (
    <form onSubmit={handleSubmit}>
      <select value={selectedOption} onChange={handleFolderChange}>
        <option value="">-- Select an option --</option>
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      <input type="text" value={url} onChange={handleInputChange} />
      <button type="submit">Send Request</button>
    </form>
  );
};


export default UrlInput;
