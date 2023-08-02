import { useEffect, useState } from 'react'
import { getMangaPreview } from '../api/api'
import { Manga, Volume, Chapter } from '../api/models'

const ChapterRaw = ({ ch, index, Prefix }: { ch: Chapter; index: string; Prefix: JSX.Element }) => {
  console.log('ch: ', index)
  return (
    <div key={index} className="flex ml-2 text-xl">
      {Prefix}
      <div className="justify-center items-center">
        <input id={index} type="checkbox" className="checkbox checkbox-xs mr-2" />
      </div>
      <p className="w-10 overflow-hidden">{ch.Info.Identifier}</p>
      <p className="ml-2">{ch.Info.Title}</p>
    </div>
  )
}

const VolumeView = ({ v, index: volumeIndex }: { v: Volume; index: string }) => {
  var chapterIDs: string[] = []
  const [checked, setChecked] = useState<string[]>([])

  function handleCheckboxesCheckAll(checked: any, isChecked: boolean) {
    setChecked(isChecked ? chapterIDs : [])
  }

  console.log('vol: ', volumeIndex)

  const midPrefix = (
    <p className="w-6" style={{ marginLeft: '3px' }}>
      ├
    </p>
  )

  const endPrefix = (
    <p className="w-6" style={{ marginRight: '3px' }}>
      └
    </p>
  )

  return (
    <div key={volumeIndex}>
      <details className="collapse  hover:bg-base-300">
        <summary className="m-2, p-2">
          <div className="flex">
            <div className="text-center w-14 flex">
              <label className="m-auto justify-center items-center">
                <input
                  id={volumeIndex}
                  // indeterminate={checked.length > 0 && checked.length < chapterIDs.length}
                  checked={checked.length === chapterIDs.length}
                  onChange={() => handleCheckboxesCheckAll}
                  type="checkbox"
                  className="checkbox"
                />
              </label>
            </div>
            <div className="flex items-center space-x-3">
              <div className="avatar">
                <div className="mask mask-squircle w-20 h-20">
                  <img src={v.CoverPath} alt="volume cover" />
                </div>
              </div>
              <div>
                <div className="font-bold text-2xl">Volume - {v.Info.Identifier}</div>
                <div className="text-sm opacity-50">{v.Chapters.length} chapters</div>
              </div>
            </div>
          </div>
        </summary>
        <div className="collapse-content">
          {v.Chapters.map((ch, chIndex) => {
            const chID = volumeIndex + '-' + chIndex.toString()
            chapterIDs.push(chID)
            if (v.Chapters.length - 1 !== chIndex) {
              return <ChapterRaw ch={ch} index={chID} Prefix={midPrefix} />
            }
            return <ChapterRaw ch={ch} index={chID} Prefix={endPrefix} />
          })}
        </div>
      </details>
      <div className="divider"></div>
    </div>
  )
}

const Upscale = () => {
  const [manga, setManga] = useState<Manga | null>(null)

  useEffect(() => {
    // declare the data fetching function
    const fetchData = async () => {
      const data = await getMangaPreview('10a4985d-0713-462e-a9d6-767bf91e4fd7')
      setManga(data)
    }

    // call the function
    fetchData()
      // make sure to catch any error
      .catch(console.error)
  }, [])

  if (manga) {
    return (
      <div className="max-w-5xl m-auto my-6" style={{ width: '50%' }}>
        <div className="overflow-x-auto">
          <div className="table">
            <div>
              {manga.Volumes.map((v, index) => (
                <VolumeView v={v} index={index.toString()} />
              ))}
            </div>
          </div>
        </div>
      </div>
    )
  }
  return <span className="loading loading-infinity loading-lg"></span>
}

export default Upscale
