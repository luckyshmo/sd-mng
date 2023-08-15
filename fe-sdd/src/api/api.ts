import axios from 'axios'
import { MangaDex, StoredManga } from './models'
import { chapterStore } from '../store/chapters'

const address = import.meta.env.VITE_API_URL!

export interface LoraInfo {
  name: string
  token: string
}

export async function getLoraInfos(): Promise<LoraInfo[]> {
  try {
    const response = await axios.get(address + 'info/lora')
    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }
    return response.data
  } catch (error: any) {
    //! how to handle this properly?
    console.error(`Error on getting loras: ${error.message}`)
    return []
  }
}

export async function getStoredMangaList(): Promise<StoredManga[]> {
  try {
    const response = await axios.get(address + 'manga/origin/info')
    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }
    return response.data
  } catch (error: any) {
    console.error(`Error on getting stored manga list: ${error.message}`)
    return []
  }
}

export async function downloadFile(url: string, folder: string): Promise<string> {
  const urlParams = new URLSearchParams({
    url: url,
    folder: folder,
  })

  try {
    const response = await axios.post(address + '?' + urlParams.toString())

    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }

    console.log(response.data)
    return ''
  } catch (error: any) {
    console.error(`Error: ${error.message}`)
    return error.message
  }
}

export async function getMangaPreview(id: string): Promise<MangaDex> {
  const urlParams = new URLSearchParams({
    id: id,
  })
  try {
    const response = await axios.get(address + 'upscale/info?' + urlParams.toString())
    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }
    return response.data
  } catch (error: any) {
    //! how to handle this properly?
    console.error(`Error on getting loras: ${error.message}`)
    throw error
  }
}

function convertMapToJSON(inputMap: Map<string, Map<string, void>>): string {
  const outputObj: Record<string, Record<string, unknown>> = {}

  inputMap.forEach((innerMap, outerKey) => {
    const innerObj: Record<string, unknown> = {}

    innerMap.forEach((_, innerKey) => {
      innerObj[innerKey] = null
    })

    outputObj[outerKey] = innerObj
  })

  return JSON.stringify(outputObj)
}

export async function downloadFilteredManga(id: string): Promise<void> {
  const chapters = chapterStore.volToChapters

  if (chapters === undefined) {
    throw new Error(`No chapters in store`)
  }

  const body = convertMapToJSON(chapters)

  const urlParams = new URLSearchParams({ id: id })
  try {
    const response = await axios.post(address + 'upscale/download?' + urlParams.toString(), body)
    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`)
    }
  } catch (error: any) {
    throw error
  }
}
