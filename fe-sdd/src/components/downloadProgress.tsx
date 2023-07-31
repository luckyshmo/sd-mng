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
          <div className="flex items-center m-1" key={index}>
            <progress
              className="progress progress-secondary w-56"
              value={item.percentage}
              max="100"
            ></progress>
            <p className="ml-2">{item.name}</p>
          </div>
        )
      })}
    </div>
  )
}

export default observer(ProgressItems)
