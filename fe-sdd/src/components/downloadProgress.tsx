import { observer } from 'mobx-react-lite'
import { downloadStore } from '../store/store'

function ProgressItems() {
  if (downloadStore.downloads.length === 0) {
    return <div>No downloads</div>
  }
  return (
    <div>
      {downloadStore.downloads.map((item, index) => {
        return (
          <div key={index}>
            <p>
              {item.name}: {item.percentage}%
            </p>
          </div>
        )
      })}
    </div>
  )
}

export default observer(ProgressItems)
