import { useEffect, useState } from 'react'
import { getMangaPreview } from '../api/api'
import { Manga, Volume, Chapter } from '../api/models'

const ChapterRaw = ({ ch }: { ch: Chapter }) => {
  return (
    <tr key={ch.Info.ID}>
      <th>
        <label>
          <input type="checkbox" className="checkbox" />
        </label>
      </th>
      <td>
        <p>{ch.Info.Identifier}</p>
      </td>
      <td>
        <p>{ch.Info.Title}</p>
      </td>
    </tr>
  )
}

const VolumeView = ({ v }: { v: Volume }) => {
  return (
    <tr className="hover">
      <details className="collapse">
        <summary>
          <th>
            <label>
              <input type="checkbox" className="checkbox" />
            </label>
          </th>
          <td>
            <div className="flex items-center space-x-3">
              <div className="avatar">
                <div className="mask mask-squircle w-12 h-12">
                  <img src={v.CoverPath} alt="Avatar Tailwind CSS Component" />
                </div>
              </div>
              <div>
                <div className="font-bold">Volume - {v.Info.Identifier}</div>
                <div className="text-sm opacity-50">{v.Chapters.length} chapters</div>
              </div>
            </div>
          </td>
        </summary>
        <div className="collapse-content">
          <table>
            <tbody>
              {v.Chapters.map((ch) => (
                <ChapterRaw ch={ch} />
              ))}
            </tbody>
          </table>
        </div>
      </details>
    </tr>
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
      <div className="max-w-5xl m-auto" style={{ width: '50%' }}>
        <div className="overflow-x-auto">
          <table className="table">
            <tbody>
              {manga.Volumes.map((v) => (
                <VolumeView v={v} />
              ))}
            </tbody>
          </table>
        </div>
      </div>
    )
  }
  return <span className="loading loading-infinity loading-lg"></span>
}

export default Upscale
