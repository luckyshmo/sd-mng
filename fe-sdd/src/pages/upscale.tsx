import { useEffect, useState } from 'react'
import { getMangaPreview } from '../api/api'
import { Manga, Volume, Chapter } from '../api/models'
import Checkbox from '../components/checkBox'
import { observer } from 'mobx-react-lite'
import { chapterStore } from '../store/chapters'

const ChapterView = ({ ch, vuid, Prefix }: { ch: Chapter; vuid: string; Prefix: JSX.Element }) => {
  const chapters = chapterStore.volToChapters.get(vuid)

  const [isChecked, setIsChecked] = useState(chapters?.has(ch.UID))

  useEffect(() => {
    const has = chapters?.has(ch.UID)
    setIsChecked(has)
  }, [chapters?.size, chapters, ch.UID])

  const handleChange = () => {
    if (!isChecked) {
      chapterStore.addChapter(vuid, ch.UID)
    } else {
      chapterStore.removeChapter(vuid, ch.UID)
    }
  }

  return (
    <div key={ch.UID} className="flex ml-2 text-xl">
      {Prefix}
      <div className="justify-center items-center">
        <input
          onChange={() => handleChange()}
          checked={isChecked}
          type="checkbox"
          className="checkbox checkbox-xs mr-2"
        />
      </div>
      <p className="w-12 overflow-hidden">{ch.Info.Identifier}</p>
      <p className="ml-2">{ch.Info.Title}</p>
    </div>
  )
}

const ChapterObserver = observer(ChapterView)

const VolumeView = ({ v }: { v: Volume }) => {
  const chapters = chapterStore.volToChapters.get(v.UID)

  const [isChecked, setIsChecked] = useState(chapters?.size === v.Chapters.length)
  const [isIndeterminate, setIsIndeterminate] = useState(
    chapters?.size !== v.Chapters.length && chapters?.size !== 0,
  )

  useEffect(() => {
    setIsChecked(chapters?.size === v.Chapters.length)
    setIsIndeterminate(chapters?.size !== v.Chapters.length && chapters?.size !== 0)
  }, [chapters?.size, v.Chapters.length])

  function handleCheckboxesCheckAll() {
    if (!isChecked) {
      chapterStore.checkVolume(v)
    } else {
      chapterStore.uncheckVolume(v.UID)
    }
    setIsChecked(!isChecked)
    setIsIndeterminate(false)
  }

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
    <div key={v.UID}>
      <details className="collapse  hover:bg-base-300">
        <summary className="m-2, p-2">
          <div className="flex">
            <div className="text-center w-14 flex">
              <label className="m-auto justify-center items-center">
                <Checkbox
                  onChange={() => {
                    handleCheckboxesCheckAll()
                  }}
                  checked={isChecked}
                  indeterminate={isIndeterminate}
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
            if (v.Chapters.length - 1 !== chIndex) {
              return <ChapterObserver ch={ch} vuid={v.UID} Prefix={midPrefix} />
            }
            return <ChapterObserver ch={ch} vuid={v.UID} Prefix={endPrefix} />
          })}
        </div>
      </details>
      <div className="divider"></div>
    </div>
  )
}

const VolumeObserver = observer(VolumeView)

const Upscale = () => {
  const [manga, setManga] = useState<Manga | null>(null)

  useEffect(() => {
    // declare the data fetching function
    const fetchData = async () => {
      const data = await getMangaPreview('10a4985d-0713-462e-a9d6-767bf91e4fd7')
      chapterStore.init(data.Volumes)
      setManga(data)
    }

    fetchData().catch(console.error)
  }, [])

  if (manga) {
    return (
      <div className="max-w-5xl m-auto my-6 overflow-x-auto table" style={{ width: '50%' }}>
        <div>
          {manga.Volumes.map((v) => (
            <VolumeObserver v={v} />
          ))}
        </div>
      </div>
    )
  }
  return <span className="loading loading-infinity loading-lg"></span>
}

export default Upscale
