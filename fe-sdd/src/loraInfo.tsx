import { useEffect, useState } from 'react'
import { getLoraInfos, LoraInfo } from './api/api'

const itemStyles = {
  display: 'flex',
  alignItems: 'center',
}

const textStyles = {
  marginRight: '10px',
}

const Item = ({ item }: { item: LoraInfo }) => {
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
    <div style={itemStyles}>
      <p style={textStyles}>
        {item.name}: {item.token}
      </p>
      <button onClick={() => copyToClipboard(item.token)}>{buttonText}</button>
    </div>
  )
}

const LoraInfoComponent = () => {
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
        <h2>LoraInfo</h2>
        <div>
          {loras.map((item, index) => {
            return <Item key={index} item={item} />
          })}
        </div>
      </div>
    )
  }

  return <div></div>
}

export default LoraInfoComponent
