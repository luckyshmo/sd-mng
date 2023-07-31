import React, { useState, ChangeEvent, FormEvent } from 'react'
import { downloadFile } from '../api/api'
import DownloadProgress from '../components/downloadProgress'

type Option = {
  value: string
  label: string
}

const options: Option[] = [
  { value: 'models/Lora', label: 'Lora' },
  { value: 'models/Stable-Diffusion', label: 'Model' },
]

const UrlInput: React.FC = () => {
  const [url, setUrl] = useState('')
  const [selectedOption, setSelectedOption] = useState(options[0].value)

  const handleFolderChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedOption(event.target.value)
  }

  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    setUrl(event.target.value)
  }

  function isValidUrl(url: string): boolean {
    // Regular expression to validate URLs
    const urlRegex = new RegExp(
      '^(https?:\\/\\/)?' + // protocol
        '((([a-zA-Z\\d]([a-zA-Z\\d-]{0,61}[a-zA-Z\\d])?)\\.)+[a-zA-Z]{2,}|' + // domain name
        '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
        '(\\:\\d+)?(\\/[-a-zA-Z\\d%@_.~+&:]*)*' + // port and path
        '(\\?[;&a-zA-Z\\d%@_.,~+&:=-]*)?' + // query string
        '(\\#[-a-zA-Z\\d_]*)?$', // fragment locator
      'i',
    )

    return urlRegex.test(url)
  }

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!isValidUrl(url)) {
      alert('Invalid URL')
      return
    }

    const message = await downloadFile(url, selectedOption)
    if (message) {
      alert(message)
    }
  }

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <select
          className="select select-secondary w-full max-w-xs"
          value={selectedOption}
          onChange={handleFolderChange}
        >
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        <input
          type="text"
          className="input input-bordered input-secondary w-full max-w-xs"
          value={url}
          onChange={handleInputChange}
        />
        <button className="btn btn-primary" type="submit">
          Send Request
        </button>
      </form>
      <DownloadProgress />
    </div>
  )
}

export default UrlInput
