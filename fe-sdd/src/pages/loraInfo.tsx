import { useEffect, useState } from 'react'
import { getLoraInfos, LoraInfo } from '../api/api'

const Item = ({ loraInfo }: { loraInfo: LoraInfo }) => {
  const defaultBtnText = 'COPY'
  const [buttonText, setButtonText] = useState(defaultBtnText)

  const copyToClipboard = (token: string) => {
    navigator.clipboard
      .writeText(token)
      .then(() => {
        console.log('Copied to clipboard:', token)
        setButtonText('✔✔✔')

        setTimeout(() => {
          setButtonText(defaultBtnText)
        }, 3000) // Reset button text after 3 seconds
      })
      .catch((error) => {
        console.error('Failed to copy to clipboard:', error)
      })
  }

  return (
    <div className="flex items-center m-1">
      <button
        className="btn btn-primary btn-xs btn-outline swap"
        onClick={() => copyToClipboard(loraInfo.token)}
      >
        {buttonText}
      </button>
      <p className="ml-2">{loraInfo.name}</p>
    </div>
  )
}

const LoraInfoComponent: React.FC = () => {
  const [loras, setLora] = useState<LoraInfo[]>([])

  const getLora = async () => {
    setLora(await getLoraInfos())
  }

  useEffect(() => {
    getLora()
  }, [])

  if (loras.length > 0) {
    return (
      <div>
        <div>
          {loras.map((item, index) => {
            return <Item key={index} loraInfo={item} />
          })}
        </div>
      </div>
    )
  }

  return <div></div>
}

export default LoraInfoComponent
