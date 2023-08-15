import { useEffect, useState } from 'react'
import { getMangaPreview, downloadFilteredManga, getStoredMangaList } from '../api/api'
import { MangaDex, StoredManga, Volume, Chapter } from '../api/models'
import Checkbox from '../components/checkBox'
import { observer } from 'mobx-react-lite'
import { chapterStore } from '../store/chapters'
import { AiOutlineDownload, AiOutlineClear } from 'react-icons/ai'

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
    <div className="flex ml-2 text-xl">
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
    <div>
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
              return <ChapterObserver key={chIndex} ch={ch} vuid={v.UID} Prefix={midPrefix} />
            }
            return <ChapterObserver key={chIndex} ch={ch} vuid={v.UID} Prefix={endPrefix} />
          })}
        </div>
      </details>
      <div className="divider"></div>
    </div>
  )
}

const VolumeObserver = observer(VolumeView)

const MangaView = ({ manga }: { manga: MangaDex }) => {
  return (
    <div>
      {manga.Volumes.map((v, ind) => (
        <VolumeObserver key={ind} v={v} />
      ))}
    </div>
  )
}

const StoredMangaView = ({ manga }: { manga: StoredManga[] }) => {
  return (
    <div>
      {manga.map((m, ind) => (
        <div key={ind} className="m-2 p-2 hover:bg-base-300 rounded-3xl">
          <div className="flex justify-between">
            <div className="flex items-center space-x-3">
              <div className="avatar">
                <div className="mask mask-squircle w-20 h-20">
                  <img src={'data:image/jpg;base64, ' + m.Cover} alt="volume cover" />
                </div>
              </div>
              <div>
                <div className="font-bold text-2xl">{m.Title}</div>
                <div className="text-sm opacity-50">{m.VolumeCount} volumes</div>
              </div>
            </div>
            <div className="flex items-center">
              <button onClick={() => {}} className="btn btn-primary">
                Upscale
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

const Upscale = () => {
  const [manga, setManga] = useState<MangaDex | null>(null)
  const [mangaID, setMangaID] = useState<string>('296cbc31-af1a-4b5b-a34b-fee2b4cad542')
  const [isLoading, setIsLoading] = useState<boolean>(false)
  const [storedManga, setStoredManga] = useState<StoredManga[] | null>(null)

  useEffect(() => {
    async function fetchData() {
      try {
        const storedManga = await getStoredMangaList()
        setStoredManga(storedManga)
      } catch (error) {
        console.error('Error fetching data:', error)
      }
    }

    fetchData()
  }, [])

  const fetchData = async () => {
    setIsLoading(true)
    const mangaPreview = await getMangaPreview(mangaID).catch(console.error)
    if (mangaPreview) {
      chapterStore.init(mangaPreview.Volumes)
      setManga(mangaPreview)
    }
    setIsLoading(false)
  }

  const handleClear = () => {
    setIsLoading(false)
    setManga(null)
    setMangaID('296cbc31-af1a-4b5b-a34b-fee2b4cad542')
  }

  const handleDownload = async () => {
    setIsLoading(true)
    await downloadFilteredManga(mangaID).catch(console.error)
    setIsLoading(false)
  }

  if (isLoading) {
    return <span className="loading loading-infinity loading-lg"></span>
  }

  if (manga) {
    return (
      <div>
        <div className="m-auto my-6 overflow-x-auto table max-w-3xl">
          {/* download and clear buttons */}
          <div className="absolute right-[-75px] flex flex-col border border-base-300 rounded-2xl">
            <div onClick={handleDownload} className="w-16 h-12 hover:bg-base-300 rounded-t-2xl">
              <AiOutlineDownload className="cursor-pointer w-8 h-8 mx-4 mb-4 mt-2" />
            </div>
            <div onClick={handleClear} className="w-16 h-12 hover:bg-base-300 rounded-b-2xl">
              <AiOutlineClear className="cursor-pointer w-8 h-8 mx-4 mb-4 mt-2" />
            </div>
          </div>
          {/* manga view */}
          <MangaView manga={manga} />
        </div>
      </div>
    )
  }

  return (
    <div className="m-auto my-4">
      <div className="flex">
        <input
          type="text"
          placeholder="Type manga id..."
          className="input input-bordered mx-4 w-96"
          value={mangaID}
          onChange={(e) => setMangaID(e.target.value)}
        />
        <button onClick={fetchData} className="btn btn-primary">
          Fetch
        </button>
      </div>
      <div className="m-auto my-6 overflow-x-auto table max-w-3xl">
        <StoredMangaView manga={storedManga || []} />
      </div>
    </div>
  )
}

export default Upscale
